const net = require('net');
const readline = require('readline');

const port = 3000;
const host = 'localhost';

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

let playerName = null;

const client = net.createConnection({ port, host }, () => {
    console.log("Connecté au serveur.");
    askForName();
});

client.setEncoding('utf8');

function sendMessage(message) {
    client.write(JSON.stringify(message) + "\n");
}

function askForName() {
    rl.question("Entrez votre nom: ", (name) => {
        playerName = name.trim() || "Joueur";
        sendMessage({ type: "join", payload: { name: playerName } });
    });
}

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

client.on('close', () => {
    console.log("Déconnecté du serveur.");
    process.exit(0);
});
