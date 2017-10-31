# cryptage
Cheating-resistant peer2peer cribbage

[![Build Status](https://travis-ci.org/whereswaldon/cryptage.svg?branch=master)](https://travis-ci.org/whereswaldon/cryptage)

NOTE: This is currently in heavy development and is likely to be unstable.

This project is an attempt to write a cheat-resistant p2p card game implementation. In a client-server implementation of a card game, a trusted server can hold the state of the deck so that neither client has it (and therefore neigher client can cheat by peeking at it). In a p2p context, both players need to hold a copy of the deck. How can you prevent either player from cheating by looking at the deck in memory?

Encrypt the deck.

That's what this project is an implementation of. The end goal is to provide:
1. A library that exposes a Deck for use implementing any given card game
2. An implementation of Cribbage using that library
