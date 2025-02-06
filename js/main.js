// main.js
const { createInterface, askQuestion } = require('./src/utils');
const { startLocalGame } = require('./src/local');
const { startServer } = require('./src/server');

async function main() {
    const rl = createInterface();
    console.log("=== Bienvenue dans Just One ===");
    console.log("Choisissez le mode de jeu :");
    console.log("1 - Jeu local (tous les joueurs sur un même terminal)");
    console.log("2 - Jeu multi-terminal");

    const mode = await askQuestion(rl, "Votre choix (1 ou 2): ");
    if (mode.trim() === "1") {
        rl.close();
        startLocalGame();
    } else if (mode.trim() === "2") {
        console.log("\nChoisissez votre rôle dans le mode multi-terminal :");
        console.log("1 - Hôte (crée la partie)");
        console.log("2 - Client (rejoint une partie existante)");
        const role = await askQuestion(rl, "Votre choix (1 ou 2): ");
        rl.close();
        if (role.trim() === "1") {
            startServer();
        } else if (role.trim() === "2") {
            console.log("Lancement du client...");
            require('./src/client'); // Charger dynamiquement le client
        } else {
            console.log("Choix invalide.");
            process.exit(1);
        }
    } else {
        console.log("Choix invalide.");
        process.exit(1);
    }
}

main();
