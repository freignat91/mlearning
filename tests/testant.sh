ml network loadTrainFile ./tests/testant.json
ml network train ant -n $1 -a -c --hide
ml network loadTrainFile ./tests/testa8.json
ml network train testa8 -n 1
