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
        // On mélange le deck de mots pour avoir un ordre aléatoire
        this.deck = this.shuffle([...words]);
        this.successful = [];
        this.discarded = [];
    }

    // Méthode de démarrage de la partie
    async start() {
        console.log("\n=== Début de la partie ===\n");

        // Boucle tant que le deck n'est pas épuisé
        while (this.deck.length > 0) {
            await this.playRound();
            // Passer le rôle de joueur actif au joueur suivant
            this.players.rotate();
        }

        // Fin de la partie, affichage du score
        console.log("\n=== Fin de la partie ===");
        console.log(`Cartes réussies : ${this.successful.length}`);
        console.log(`Cartes ratées : ${this.discarded.length}`);
        this.rl.close();
    }

    // Méthode pour mélanger un tableau
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

    // Méthode pour gérer un tour de jeu
    async playRound() {
        console.log(`\n--- Tour ${this.currentRound + 1} ---`);
        if (this.deck.length === 0) return;

        // Tire le mot mystère pour ce tour
        const currentWord = this.deck.shift();
        // Dans un jeu en vrai, le mot mystère ne doit être vu que par le joueur actif.
        // Ici, on affiche le mot pour simplifier la démonstration.
        console.log(`(Pour le joueur actif uniquement) Mot mystère : ${currentWord}`);

        const activePlayer = this.players.getActivePlayer();
        console.log(`Le joueur actif est : ${activePlayer}`);

        // Collecte des indices fournis par les autres joueurs
        let indices = {};
        for (let player of this.players.list) {
            if (player === activePlayer) continue; // Le joueur actif ne donne pas d'indice
            const answer = await this.askPlayerForIndex(player);
            indices[player] = answer.trim().toLowerCase();
        }

        // Filtrage des indices identiques : annulation des indices qui apparaissent plus d'une fois
        const validIndices = this.filterIndices(indices);

        console.log("\nIndices validés (ceux qui ne sont pas annulés) :");
        for (let key in validIndices) {
            console.log(`${key}: ${validIndices[key]}`);
        }

        // Le joueur actif essaie de deviner le mot mystère à partir des indices validés
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

    // Demande à un joueur de fournir son indice
    askPlayerForIndex(player) {
        return new Promise(resolve => {
            this.rl.question(`Joueur ${player}, donne ton indice : `, answer => {
                resolve(answer);
            });
        });
    }

    // Demande au joueur actif de deviner le mot mystère
    askGuess(activePlayer) {
        return new Promise(resolve => {
            this.rl.question(`\n${activePlayer}, fais ta proposition : `, answer => {
                resolve(answer);
            });
        });
    }

    // Filtre les indices pour éliminer les doublons
    filterIndices(indices) {
        const frequency = {};
        // Calculer la fréquence de chaque indice
        for (let key in indices) {
            const word = indices[key];
            frequency[word] = (frequency[word] || 0) + 1;
        }

        // Seuls les indices apparaissant une fois sont considérés comme valides
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
