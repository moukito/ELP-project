/**
 * @file client.js
 * @description Gère le client du jeu Just One en mode multi-terminal.
 * Ce module permet à un joueur de se connecter au serveur, d’envoyer et de recevoir des messages.
 * Il affiche les informations nécessaires à chaque joueur et interagit avec le serveur via TCP.
 */

const net = require('net');
const readline = require('readline');

const port = 3000;
const host = 'localhost';

/**
 * Interface readline pour gérer l’entrée utilisateur dans le terminal.
 */
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

let playerName = null;

/**
 * Initialise la connexion au serveur.
 * Dès la connexion, il demande au joueur de renseigner son nom.
 */
const client = net.createConnection({ port, host }, () => {
    console.log("Connecté au serveur.");
    askForName();
});

client.setEncoding('utf8');

/**
 * Envoie un message au serveur sous forme de JSON.
 * @param {Object} message - Objet contenant le type de message et les données associées.
 */
function sendMessage(message) {
    client.write(JSON.stringify(message) + "\n");
}

/**
 * Demande au joueur de saisir son nom et l'envoie au serveur.
 */
function askForName() {
    rl.question("Entrez votre nom: ", (name) => {
        playerName = name.trim() || "Joueur";
        sendMessage({ type: "join", payload: { name: playerName } });
    });
}

/**
 * Gère la réception des messages du serveur et déclenche les actions appropriées.
 */
client.on('data', (data) => {
    data.split('\n').forEach(raw => {
        if (!raw.trim()) return;
        try {
            const message = JSON.parse(raw);
            handleMessage(message);
        } catch (e) {
            console.error("Erreur de parsing :", raw);
        }
    });
});

/**
 * Analyse les messages reçus du serveur et agit en conséquence.
 * @param {Object} message - Objet contenant le type du message et ses données.
 */
function handleMessage(message) {
    switch (message.type) {
        case "join_ack":
            console.log(message.payload.message);
            break;
        case "game_start":
            console.log(message.payload.message);
            break;
        case "mystery_word":
            console.log(`Mot mystère: ${message.payload.word}`);
            break;
        case "active_notice":
            console.log(message.payload.message);
            break;
        case "ask_index":
            rl.question("Donne ton indice: ", (index) => {
                sendMessage({ type: "index", payload: { name: playerName, index } });
            });
            break;
        case "indices":
            console.log("Indices validés par tes coéquipiers:");
            console.log(message.payload.validIndices);
            rl.question("Fais ta proposition: ", (guess) => {
                sendMessage({ type: "guess", payload: { name: playerName, guess } });
            });
            break;
        case "round_result":
            if (message.payload.result === "success") {
                console.log(`Bonne réponse ! Le mot était: ${message.payload.word}`);
            } else {
                console.log(`Mauvaise réponse. Le mot était: ${message.payload.word}`);
            }
            break;
        case "game_over":
            console.log("Fin de la partie.");
            console.log(`Cartes réussies: ${message.payload.successful}`);
            console.log(`Cartes ratées: ${message.payload.discarded}`);
            process.exit(0);
            break;
        default:
            console.log("Message inconnu:", message);
    }
}

/**
 * Gère la déconnexion du client.
 */
client.on('close', () => {
    console.log("Déconnecté du serveur.");
    process.exit(0);
});
