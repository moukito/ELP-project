const readline = require('readline');
const Game = require('./game');
const Players = require('./players');
const Words = require('./words');

// Création de l'interface readline pour interagir dans le terminal
const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

// Fonction utilitaire pour poser une question et attendre la réponse
function askQuestion(query) {
    return new Promise(resolve => {
        rl.question(query, answer => {
            resolve(answer);
        });
    });
}

// Fonction asynchrone pour configurer et lancer le jeu
async function setupGame() {
    console.log("Bienvenue dans Just One (version terminal) !");

    const numPlayersStr = await askQuestion("Combien de joueurs ? ");
    const numPlayers = parseInt(numPlayersStr);
    if (isNaN(numPlayers) || numPlayers < 2) {
        console.log("Le jeu nécessite au moins 2 joueurs.");
        rl.close();
        return;
    }

    let playersList = [];
    for (let i = 0; i < numPlayers; i++) {
        const name = await askQuestion(`Nom du joueur ${i + 1} : `);
        playersList.push(name.trim() || `Joueur${i + 1}`);
    }

    // Créer la liste des joueurs
    const players = new Players(playersList);

    // Charger les mots depuis le fichier JSON
    try {
        const words = await Words.loadWords();
        // Créer et démarrer le jeu
        const game = new Game(players, words, rl);
        game.start();
    } catch (err) {
        console.error("Erreur lors du chargement des mots :", err);
        rl.close();
    }
}

setupGame();
