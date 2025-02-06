/**
 * @file game.js
 * @description Gère la logique principale du jeu Just One en mode local.
 * Ce module contrôle le déroulement de la partie, le tirage des mots mystères,
 * la gestion des indices et la validation des propositions des joueurs.
 */

class Game {
    /**
     * Initialise une nouvelle partie de Just One.
     * @constructor
     * @param {Players} players - Instance de Players qui gère la liste des joueurs.
     * @param {Array<string>} words - Tableau contenant les mots mystères.
     * @param {Interface} rl - L'interface readline pour gérer les entrées/sorties du terminal.
     */
    constructor(players, words, rl) {
        this.players = players;
        this.words = words;
        this.rl = rl;
        this.currentRound = 0;
        this.deck = this.shuffle([...words]); // Mélange les mots pour créer une pioche aléatoire
        this.successful = [];
        this.discarded = [];
    }

    /**
     * Démarre la partie et enchaîne les tours jusqu'à épuisement des mots mystères.
     * @async
     */
    async start() {
        console.log("\n=== Début de la partie ===\n");

        while (this.deck.length > 0) {
            await this.playRound();
            this.players.rotate();
        }

        console.log("\n=== Fin de la partie ===");
        console.log(`Cartes réussies : ${this.successful.length}`);
        console.log(`Cartes ratées : ${this.discarded.length}`);
        this.rl.close();
    }

    /**
     * Mélange un tableau de manière aléatoire (algorithme de Fisher-Yates).
     * @param {Array} array - Tableau à mélanger.
     * @returns {Array} - Tableau mélangé.
     */
    shuffle(array) {
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
     * Gère un tour de jeu :
     * - Tire un mot mystère
     * - Collecte et valide les indices des joueurs passifs
     * - Demande au joueur actif de deviner le mot
     * - Met à jour les scores
     * @async
     */
    async playRound() {
        console.log(`\n--- Tour ${this.currentRound + 1} ---`);
        if (this.deck.length === 0) return;

        const currentWord = this.deck.shift();
        console.log(`(Pour le joueur actif uniquement) Mot mystère : ${currentWord}`);

        const activePlayer = this.players.getActivePlayer();
        console.log(`Le joueur actif est : ${activePlayer}`);

        let indices = {};
        for (let player of this.players.list) {
            if (player === activePlayer) continue;
            const answer = await this.askPlayerForIndex(player);
            indices[player] = answer.trim().toLowerCase();
        }

        // Filtrage des indices identiques ou invalides
        const validIndices = this.filterIndices(indices);

        console.log("\nIndices validés (ceux qui ne sont pas annulés) :");
        for (let key in validIndices) {
            console.log(`${key}: ${validIndices[key]}`);
        }

        // Demande au joueur actif de deviner le mot
        const guess = await this.askGuess(activePlayer);

        // Vérifie la réponse et met à jour les scores
        if (guess.trim().toLowerCase() === currentWord.trim().toLowerCase()) {
            console.log("Bonne réponse !");
            this.successful.push(currentWord);
        } else {
            console.log("Mauvaise réponse.");
            this.discarded.push(currentWord);
        }
        this.currentRound++;
    }

    /**
     * Demande à un joueur passif de fournir un indice.
     * @param {string} player - Nom du joueur passif.
     * @returns {Promise<string>} - Indice donné par le joueur.
     * @async
     */
    askPlayerForIndex(player) {
        return new Promise(resolve => {
            this.rl.question(`Joueur ${player}, donne ton indice : `, answer => {
                resolve(answer);
            });
        });
    }

    /**
     * Demande au joueur actif de deviner le mot mystère.
     * @param {string} activePlayer - Nom du joueur actif.
     * @returns {Promise<string>} - Proposition du joueur actif.
     * @async
     */
    askGuess(activePlayer) {
        return new Promise(resolve => {
            this.rl.question(`\n${activePlayer}, fais ta proposition : `, answer => {
                resolve(answer);
            });
        });
    }

    /**
     * Filtre les indices pour supprimer ceux qui sont identiques ou invalides.
     * @param {Object} indices - Objet contenant les indices proposés par les joueurs passifs.
     * @returns {Object} - Indices validés après filtrage.
     */
    filterIndices(indices) {
        const frequency = {};
        for (let key in indices) {
            const word = indices[key];
            frequency[word] = (frequency[word] || 0) + 1;
        }

        const valid = {};
        for (let key in indices) {
            const word = indices[key];
            if (frequency[word] === 1) {
                valid[key] = word;
            }
        }
        return valid;
    }
}

module.exports = Game;
