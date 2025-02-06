class Game {
    /**
     * @param {Players} players - Instance de Players qui gère la liste des joueurs.
     * @param {Array} words - Tableau contenant les mots mystères.
     * @param {Interface} rl - L'interface readline pour les entrées/sorties.
     */
    constructor(players, words, rl) {
        this.players = players;
        this.words = words; // Tableau de mots
        this.rl = rl;
        this.currentRound = 0;
        this.deck = this.shuffle([...words]);
        this.successful = [];
        this.discarded = [];
    }

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

        const validIndices = this.filterIndices(indices);

        console.log("\nIndices validés (ceux qui ne sont pas annulés) :");
        for (let key in validIndices) {
            console.log(`${key}: ${validIndices[key]}`);
        }

        const guess = await this.askGuess(activePlayer);

        if (guess.trim().toLowerCase() === currentWord.trim().toLowerCase()) {
            console.log("Bonne réponse !");
            this.successful.push(currentWord);
        } else {
            console.log("Mauvaise réponse.");
            this.discarded.push(currentWord);
        }
        this.currentRound++;
    }

    askPlayerForIndex(player) {
        return new Promise(resolve => {
            this.rl.question(`Joueur ${player}, donne ton indice : `, answer => {
                resolve(answer);
            });
        });
    }

    askGuess(activePlayer) {
        return new Promise(resolve => {
            this.rl.question(`\n${activePlayer}, fais ta proposition : `, answer => {
                resolve(answer);
            });
        });
    }

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
