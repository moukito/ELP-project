// src/server.js
const net = require('net');
const { createInterface, askQuestion } = require('./utils');
const Words = require('./words');

const rl = createInterface();
const port = 3000;

let totalPlayers = 0;
const players = []; // Chaque joueur : { name, socket, isHost }
let hostPlayer = null;

let wordsDeck = [];
let currentRound = 0;
let successful = [];
let discarded = [];
let currentRoundWord = null;
let roundResponses = {};

// Envoie un message au joueur (si isHost, affiche sur la console)
function sendToPlayer(player, message) {
    if (player.isHost) {
        console.log(`\n[Message pour ${player.name} - HOST] ${message.type}:`, message.payload);
    } else {
        player.socket.write(JSON.stringify(message) + "\n");
    }
}

function broadcast(message) {
    players.forEach(player => sendToPlayer(player, message));
}

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
            askQuestion(rl, `Ton indice (passif) : `).then(answer => {
                roundResponses[player.name] = answer.trim().toLowerCase();
                checkIndices(activePlayer);
            });
        } else {
            sendToPlayer(player, { type: "ask_index", payload: { message: "Envoie ton indice pour ce tour." } });
        }
    });
}

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
        askQuestion(rl, `Fais ta proposition: `).then(answer => {
            processGuess(activePlayer, answer.trim().toLowerCase());
        });
    }
}

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

async function setupServer() {
    const numPlayersStr = await askQuestion(rl, "Nombre total de joueurs (y compris toi) : ");
    totalPlayers = parseInt(numPlayersStr);
    if (isNaN(totalPlayers) || totalPlayers < 2) {
        console.log("Il faut au moins 2 joueurs.");
        process.exit(1);
    }
    const hostName = await askQuestion(rl, "Ton nom (host) : ");
    hostPlayer = { name: hostName.trim() || "Host", isHost: true };
    players.push(hostPlayer);
    console.log(`En attente de ${totalPlayers - 1} joueurs supplémentaires...`);

    const words = await Words.loadWords();
    wordsDeck = shuffle(words);

    server.listen(port, () => {
        console.log(`Serveur lancé sur le port ${port}.`);
    });
}

const server = net.createServer();

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
        const playerName = message.payload.name;
        const indexText = message.payload.index;
        roundResponses[playerName] = indexText.trim().toLowerCase();
        const activePlayer = players[currentRound % players.length];
        checkIndices(activePlayer);
    } else if (message.type === "guess") {
        const playerName = message.payload.name;
        const guess = message.payload.guess.trim().toLowerCase();
        processGuess({ name: playerName, isHost: false }, guess);
    }
}

async function startServer() {
    await setupServer();
}

module.exports = { startServer };
