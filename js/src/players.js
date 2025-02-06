/**
 * @file players.js
 * @description Gère la liste des joueurs et la rotation du joueur actif.
 * Ce module permet d'initialiser les joueurs, d'obtenir le joueur actif
 * et de faire tourner l'ordre des tours.
 */

class Players {
    /**
     * Initialise la liste des joueurs et définit le premier joueur actif.
     * @constructor
     * @param {Array<string>} names - Tableau contenant les noms des joueurs.
     */
    constructor(names) {
        this.list = names;
        this.currentIndex = 0; // Index du joueur actif
    }

    /**
     * Retourne le nom du joueur actuellement actif.
     * @returns {string} - Nom du joueur actif.
     */
    getActivePlayer() {
        return this.list[this.currentIndex];
    }

    /**
     * Passe au joueur suivant dans la liste en utilisant une rotation circulaire.
     */
    rotate() {
        this.currentIndex = (this.currentIndex + 1) % this.list.length;
    }
}

module.exports = Players;
