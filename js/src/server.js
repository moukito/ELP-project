/**
 * @file server.js
 * @description Gère le serveur du jeu Just One en mode multi-terminal.
 * Ce module crée un serveur TCP permettant aux joueurs de se connecter, de recevoir des mots mystères,
 * d'envoyer des indices et de deviner le mot. Le serveur gère la logique du jeu et la communication entre les joueurs.
 */

const net = require('net');
const readline = require('readline');
const Words = require('./words');

/**
 * Mélange un tableau de manière aléatoire (algorithme de Fisher-Yates).
 * @param {Array} array - Tableau à mélanger.
 * @returns {Array} - Tableau mélangé.
 */
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

/**
 * Crée une interface readline pour gérer l'entrée utilisateur sur le terminal du serveur.
 */
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

/**
 * Pose une question à l'utilisateur et retourne sa réponse.
 * @param {string} query - Texte de la question.
 * @returns {Promise<string>} - Réponse de l'utilisateur.
 */
function askQuestion(query) {
    return new Promise(resolve => {
        rl.question(query, resolve);
    });
}

// Liste des joueurs connectés
const players = [];
let hostPlayer = null;

let totalPlayers = 0;
const port = 3000;

let wordsDeck = [];
let currentRound = 0;
let successful = [];
let discarded = [];
let currentRoundWord = null;
let roundResponses = {};

/**
 * Envoie un message à un joueur spécifique.
 * @param {Object} player - Objet représentant le joueur ({ name, socket, isHost }).
 * @param {Object} message - Objet JSON contenant le type et les données du message.
 */
function sendToPlayer(player, message) {
    if (player.isHost) {
        console.log(`\n[Message pour ${player.name} - HOST] ${message.type}:`, message.payload);
    } else {
        player.socket.write(JSON.stringify(message) + "\n");
    }
}

/**
 * Diffuse un message à tous les joueurs connectés.
 * @param {Object} message - Objet JSON contenant le type et les données du message.
 */
function broadcast(message) {
    players.forEach(player => sendToPlayer(player, message));
}

/**
 * Démarre un tour de jeu, choisit un joueur actif, distribue le mot mystère aux autres joueurs
 * et demande les indices.
 */
async function playRound() {
    if (wordsDeck.length === 0) {
        console.log("Plus de mots dans le deck. Fin de la partie.");
        broadcast({ type: "game_over", payload: { successful: successful.length, discarded: discarded.length } });
        process.exit(0);
    }
    console.log(`\n--- Tour ${currentRound + 1} ---`);

    const activePlayer = players[currentRound % players.length];
    console.log(`Joueur actif: ${activePlayer.name}`);

    currentRoundWord = wordsDeck.shift();

    players.forEach(player => {
        if (player === activePlayer) {
            sendToPlayer(player, { type: "active_notice", payload: { message: "Tu es le joueur actif. Tu ne vois pas le mot mystère." } });
        } else {
            sendToPlayer(player, { type: "mystery_word", payload: { word: currentRoundWord } });
        }
    });

    roundResponses = {};

    const passivePlayers = players.filter(p => p !== activePlayer);
    passivePlayers.forEach(player => {
        if (player.isHost) {
            askQuestion(`Ton indice (passif) : `).then(answer => {
                roundResponses[player.name] = answer.trim().toLowerCase();
                checkIndices(activePlayer);
            });
        } else {
            sendToPlayer(player, { type: "ask_index", payload: { message: "Envoie ton indice pour ce tour." } });
        }
    });
}

/**
 * Vérifie les indices fournis par les joueurs passifs, élimine les doublons et envoie les indices validés au joueur actif.
 * @param {Object} activePlayer - Joueur actif du tour.
 */
function checkIndices(activePlayer) {
    const expectedCount = players.filter(p => p !== activePlayer).length;
    if (Object.keys(roundResponses).length < expectedCount) return;

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

    sendToPlayer(activePlayer, { type: "indices", payload: { validIndices } });

    if (activePlayer.isHost) {
        askQuestion(`Fais ta proposition: `).then(answer => {
            processGuess(activePlayer, answer.trim().toLowerCase());
        });
    }
}

/**
 * Vérifie la proposition du joueur actif et met à jour les scores.
 * @param {Object} activePlayer - Joueur actif qui a fait une proposition.
 * @param {string} guess - Proposition du joueur.
 */
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
    setTimeout(playRound, 1000);
}

/**
 * Configure le serveur et attend la connexion des joueurs.
 */
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

    const words = await Words.loadWords();
    wordsDeck = shuffle(words);

    server.listen(port, () => {
        console.log(`Serveur lancé sur le port ${port}.`);
    });
}

// Création du serveur TCP
const server = net.createServer();

/**
 * Gère la connexion d'un nouveau joueur.
 */
server.on('connection', (socket) => {
    console.log("Un joueur distant vient de se connecter.");
    socket.setEncoding('utf8');

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
    });
});

/**
 * Gère les messages reçus des clients et effectue les actions appropriées.
 * @param {Object} socket - Socket du client.
 * @param {Object} message - Message JSON reçu.
 */
function handleClientMessage(socket, message) {
    if (message.type === "join") {
        const playerName = message.payload.name;
        const newPlayer = { name: playerName, socket, isHost: false };
        players.push(newPlayer);
        console.log(`Le joueur ${playerName} a rejoint la partie.`);
        socket.write(JSON.stringify({ type: "join_ack", payload: { message: "Bienvenue dans la partie!" } }) + "\n");

        if (players.length === totalPlayers) {
            console.log("Tous les joueurs sont connectés. La partie va commencer !");
            broadcast({ type: "game_start", payload: { message: "La partie commence !" } });
            playRound();
        }
    } else if (message.type === "index") {
        roundResponses[message.payload.name] = message.payload.index.trim().toLowerCase();
        checkIndices(players[currentRound % players.length]);
    } else if (message.type === "guess") {
        processGuess(players[currentRound % players.length], message.payload.guess.trim().toLowerCase());
    }
}

// Initialisation du serveur
(async function init() {
    await setupServer();
})();
