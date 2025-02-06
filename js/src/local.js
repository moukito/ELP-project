// src/local.js
const Players = require('./players');
const Game = require('./game');
const Words = require('./words');
const { createInterface, askQuestion } = require('./utils');

async function startLocalGame() {
    const rl = createInterface();
    console.log("Bienvenue dans Just One (version locale) !");

    const numPlayersStr = await askQuestion(rl, "Combien de joueurs ? ");
    const numPlayers = parseInt(numPlayersStr);
    if (isNaN(numPlayers) || numPlayers < 2) {
        console.log("Le jeu nécessite au moins 2 joueurs.");
        rl.close();
        return;
    }

    let playersList = [];
    for (let i = 0; i < numPlayers; i++) {
        const name = await askQuestion(rl, `Nom du joueur ${i + 1} : `);
        playersList.push(name.trim() || `Joueur${i + 1}`);
    }

    const players = new Players(playersList);

    try {
        const words = await Words.loadWords();
        const game = new Game(players, words, rl);
        game.start();
    } catch (err) {
        console.error("Erreur lors du chargement des mots :", err);
        rl.close();
    }
}

module.exports = { startLocalGame };
