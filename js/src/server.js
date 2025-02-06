const net = require('net');
const readline = require('readline');
const Words = require('./words');

// Utilitaire pour mélanger un tableau
function shuffle(array) {
    let currentIndex = array.length, temporaryValue, randomIndex;
    while (currentIndex !== 0) {
        randomIndex = Math.floor(Math.random() * currentIndex);
        currentIndex -= 1;
        temporaryValue = array[currentIndex];
        array[currentIndex] = array[randomIndex];
        array[randomIndex] = temporaryValue;
    }
    return array;
}

// Interface readline pour l'hôte (serveur local)
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

function askQuestion(query) {
    return new Promise(resolve => {
        rl.question(query, resolve);
    });
}

// Liste des joueurs : chaque joueur est un objet { name, socket, isHost }
// Pour l'hôte, aucune socket n'est utilisée (les échanges se font via la console)
const players = [];
let hostPlayer = null;

// Configuration de la partie
let totalPlayers = 0;
const port = 3000;

// Variables de jeu
let wordsDeck = [];
let currentRound = 0;
let successful = [];
let discarded = [];
let currentRoundWord = null; // mot mystère courant
let roundResponses = {};   // Réponses (indices) recueillies pour le tour

// Fonctions d'envoi de messages
function sendToPlayer(player, message) {
    if (player.isHost) {
        // Pour l'hôte, on affiche le message dans la console (précédé du nom)
        console.log(`\n[Message pour ${player.name} - HOST] ${message.type}:`, message.payload);
    } else {
        // Pour un joueur distant, on envoie le JSON suivi d'un saut de ligne
        player.socket.write(JSON.stringify(message) + "\n");
    }
}

function broadcast(message) {
    players.forEach(player => sendToPlayer(player, message));
}

// Gestion du tour de jeu
async function playRound() {
    if (wordsDeck.length === 0) {
        console.log("Plus de mots dans le deck. Fin de la partie.");
        broadcast({ type: "game_over", payload: { successful: successful.length, discarded: discarded.length } });
        process.exit(0);
    }
    console.log(`\n--- Tour ${currentRound + 1} ---`);

    // Détermine le joueur actif (rotation circulaire)
    const activePlayer = players[currentRound % players.length];
    console.log(`Joueur actif: ${activePlayer.name}`);

    // Tire le mot mystère et sauvegarde-le
    currentRoundWord = wordsDeck.shift();

    // Envoie aux joueurs :
    players.forEach(player => {
        if (player === activePlayer) {
            // Le joueur actif ne reçoit pas le mot mystère
            sendToPlayer(player, { type: "active_notice", payload: { message: "Tu es le joueur actif. Tu ne vois pas le mot mystère." } });
        } else {
            // Les joueurs passifs reçoivent le mot mystère
            sendToPlayer(player, { type: "mystery_word", payload: { word: currentRoundWord } });
        }
    });

    // Réinitialise les réponses pour ce tour
    roundResponses = {};

    // Demande aux joueurs passifs leur indice.
    // Pour l'hôte (s'il n'est pas actif), on le fait via la console.
    const passivePlayers = players.filter(p => p !== activePlayer);
    passivePlayers.forEach(player => {
        if (player.isHost) {
            askQuestion(`Ton indice (passif) : `).then(answer => {
                roundResponses[player.name] = answer.trim().toLowerCase();
                checkIndices(activePlayer);
            });
        } else {
            // Pour les clients distants, on envoie une requête pour qu'ils saisissent leur indice
            sendToPlayer(player, { type: "ask_index", payload: { message: "Envoie ton indice pour ce tour." } });
        }
    });
}

// Vérifie si tous les indices ont été reçus et, le cas échéant, les filtre puis les envoie au joueur actif.
function checkIndices(activePlayer) {
    const expectedCount = players.filter(p => p !== activePlayer).length;
    if (Object.keys(roundResponses).length < expectedCount) return;

    // Filtrage des indices en éliminant les doublons
    const freq = {};
    for (let key in roundResponses) {
        const idx = roundResponses[key];
        freq[idx] = (freq[idx] || 0) + 1;
    }
    const validIndices = {};
    for (let key in roundResponses) {
        const idx = roundResponses[key];
        if (freq[idx] === 1) validIndices[key] = idx;
    }

    // Envoie les indices validés au joueur actif
    sendToPlayer(activePlayer, { type: "indices", payload: { validIndices } });

    // Si le joueur actif est l'hôte, on demande directement via la console sa proposition.
    if (activePlayer.isHost) {
        askQuestion(`Fais ta proposition: `).then(answer => {
            processGuess(activePlayer, answer.trim().toLowerCase());
        });
    }
}

// Traitement de la proposition du joueur actif
function processGuess(activePlayer, guess) {
    if (guess === currentRoundWord.trim().toLowerCase()) {
        console.log("Bonne réponse !");
        successful.push(currentRoundWord);
        broadcast({ type: "round_result", payload: { result: "success", word: currentRoundWord } });
    } else {
        console.log("Mauvaise réponse.");
        discarded.push(currentRoundWord);
        broadcast({ type: "round_result", payload: { result: "fail", word: currentRoundWord } });
    }
    currentRound++;
    // Lance le tour suivant après une courte pause
    setTimeout(playRound, 1000);
}

// Configuration initiale du serveur (pour l'hôte)
async function setupServer() {
    const numPlayersStr = await askQuestion("Nombre total de joueurs (y compris toi) : ");
    totalPlayers = parseInt(numPlayersStr);
    if (isNaN(totalPlayers) || totalPlayers < 2) {
        console.log("Il faut au moins 2 joueurs.");
        process.exit(1);
    }
    const hostName = await askQuestion("Ton nom (host) : ");
    hostPlayer = { name: hostName.trim() || "Host", isHost: true };
    players.push(hostPlayer);
    console.log(`En attente de ${totalPlayers - 1} joueurs supplémentaires...`);

    // Charge les mots et mélange le deck
    const words = await Words.loadWords();
    wordsDeck = shuffle(words);

    // Démarre le serveur TCP
    server.listen(port, () => {
        console.log(`Serveur lancé sur le port ${port}.`);
    });
}

// Création du serveur TCP pour accepter les connexions clients
const server = net.createServer();

server.on('connection', (socket) => {
    console.log("Un joueur distant vient de se connecter.");
    socket.setEncoding('utf8');

    // Lorsque le serveur reçoit des données depuis un client
    socket.on('data', (data) => {
        data.split('\n').forEach(raw => {
            if (!raw.trim()) return;
            try {
                const message = JSON.parse(raw);
                handleClientMessage(socket, message);
            } catch (e) {
                console.error("Erreur de parsing du message :", raw);
            }
        });
    });

    socket.on('close', () => {
        console.log("Un joueur s'est déconnecté.");
        // (Ici, vous pourriez gérer la déconnexion d'un joueur.)
    });
});

// Gère les messages reçus des clients
function handleClientMessage(socket, message) {
    if (message.type === "join") {
        const playerName = message.payload.name;
        const newPlayer = { name: playerName, socket, isHost: false };
        players.push(newPlayer);
        console.log(`Le joueur ${playerName} a rejoint la partie.`);
        // Accusé de réception
        socket.write(JSON.stringify({ type: "join_ack", payload: { message: "Bienvenue dans la partie!" } }) + "\n");

        // Lorsque tous les joueurs sont connectés, la partie démarre.
        if (players.length === totalPlayers) {
            console.log("Tous les joueurs sont connectés. La partie va commencer !");
            broadcast({ type: "game_start", payload: { message: "La partie commence !" } });
            playRound();
        }
    } else if (message.type === "index") {
        // Réception d'un indice de la part d'un joueur distant
        const playerName = message.payload.name;
        const indexText = message.payload.index;
        roundResponses[playerName] = indexText.trim().toLowerCase();
        const activePlayer = players[currentRound % players.length];
        checkIndices(activePlayer);
    } else if (message.type === "guess") {
        // Réception de la proposition du joueur actif distant
        const playerName = message.payload.name;
        const guess = message.payload.guess.trim().toLowerCase();
        processGuess({ name: playerName, isHost: false }, guess);
    }
}

// Lancement de la configuration initiale
(async function init() {
    await setupServer();
})();
