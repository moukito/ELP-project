class Players {
    /**
     * @param {Array} names - Tableau contenant les noms des joueurs.
     */
    constructor(names) {
        this.list = names;
        this.currentIndex = 0;
    }

    getActivePlayer() {
        return this.list[this.currentIndex];
    }

    rotate() {
        this.currentIndex = (this.currentIndex + 1) % this.list.length;
    }
}

module.exports = Players;
