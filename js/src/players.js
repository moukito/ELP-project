class Players {
    /**
     * @param {Array} names - Tableau contenant les noms des joueurs.
     */
    constructor(names) {
        this.list = names;
        this.currentIndex = 0; // Indice du joueur actif
    }

    // Retourne le joueur actif
    getActivePlayer() {
        return this.list[this.currentIndex];
    }

    // Passe le tour au joueur suivant
    rotate() {
        this.currentIndex = (this.currentIndex + 1) % this.list.length;
    }
}

module.exports = Players;
