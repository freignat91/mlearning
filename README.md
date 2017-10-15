# Ants Machine learning

mlearning v0.0.3


# project purpose


Use multi-layer neuron network to drive ant behaviours, each ant having its own network.

Display graphical representation of the ant behaviours and information to monitor how they evolve and how their network converge.

Use the virtual game context to train the networks, train the networks permanently using what the ants see and the consequences of theirs actions.

version 0.0.1: Ants are able to spread the virtual space in order to cover all the available space without having contact between them. The ants avoid each other

version 0.0.2: Foods appears in the space, ants are able to get them and bring back them to the nest, letting a pheromone path and other ants are able to trace back the a pheromone path to find food sources

version 0.0.3: with two nests and two ant types (worker and soldier), ants will be able to fight agains the other nest ants.

version 0.0.4 (not done): in version 0.0.2 worker panic mode and worker return to nest are hard-coded, in 0.0.4, panic mode should be network train result.

version 0.0.5 (not done): have return to nests as network train result, including avoid obstacles on the path


# Version 0.0.3 rules of the game

In version 0.0.3 the nests try to survive maximizing their resources.

The ant behaviours are the following:

- by default, workers disperse them-self covering the maximum of space as quickly as possible
- if a worker finds a food, it brings back it to nest
- when a worker reach its nest with a food, the nest receive 1 resource point and create a new worker
- if a worker encounter an hostile ant (not belonging to the same nest), then it enter in panic mode and go back quickly to its nest
- when a worker in panic mode reach its nest, the nest create a new soldier running in the direction  of the panic source
- when a soldier encounter an hostile ant (worker or soldier), it attacks it. each contact makes by a soldier remove one life point to its target (a soldier has 400 live point, a worker 120, see ./nests/parameters.go file)

That's all, at the end, only one nest survive.


# Way of working

Ant network are train using regular gradient retro-propagation algorithm. The point is to be able to define what is the right output of a network for a given entry.

For that let's consider the consequence of an ant decision, using the following algorithm:

- for each simulation tick and for each ant:
  - build the network entry using what the ant see (see chapter `ant network structure`)
  - propagate the entry in the network
  - considering the network output, take a decision about the direction the ant should move (see chapter `take decision on output`). (note that there is no random on this part, ants move always on their network output no matter what it is)
  - move the ant to the chosen direction
  - compute the happiness of the ant (see chapter `the happiness of a ant`)
  - if the happiness is the same than the one of previous loop, do nothing and keep the direction as it is.
  - if this happiness is greater than the one of previous loop, then consider that the decision is good and train the network using the previous entry and the output corresponding the chosen direction.

This loop appeared not to be enough in order to make the network converging to stable results. It needed to not only reinforce the good decisions, but to fade the bad ones.

Then it's possible to add the following step:

- if the computed happiness is lower than the previous one, train the network with the previous entry and a fading output computed using the chosen direction (see chapter `fade decision`)

With this additional rule, the network started to converge well, but a new issue appeared:

Because of the fading, the networks results become more and more poor, meaning, the number of the distinct possible output decisions lower. For instance a network shows only 3 directions no matter the entries, where it should be able to show 8.

At the end, all the ants finished to move to only 2 directions, even if they reached the assigned goal of the version 0.0.1: they cover all the virtual space and the contacts number lower to 0, it's not good for the next version purpose.

Hopefully, it's possible to enrich the network decision capability to counter-balance the fade effect. To do that, let's use the following algorithm:

On regular basis (every 1000 or 10000 ticks), for each ant:
- compute the statistical network output distribution: for each possible decision, the number of time the decision has been taken in the period
- when decision has to be taken using network output:
  - compute the regular decision as usual (see chapter `take decision on output`)
  - look at the immediate other directions (+1 and -1 regarding the regular chosen one) and if one of them, hasn't been taken more than a given time (100 for now), then this decisions instead of the regular one.
  - This decision will be as usual reinforced or fade considering its consequence on the next loop.

This way, we train networks on decisions it doesn't show too much on its regular way. If the "forced" decisions appear to be good, they will be naturally reinforced and then their statistical distribution will be better (>100) and they become regular decisions.

The networks are now able to converge and use about all their decision capability.

Now, there is another issue: We don't know which network structure(s) should be used. 3 layers, 4 layers? how many neurones by layer?

The first tested networks were empirically set at the beginning, showing that a network 8-7-8 is not bad for the purpose of the version 0.0.1 (see chapter `ant network structure`), but perhaps 8-30-100-8 is better, especially for the next versions purposes...

Then, it's possible to add a mechanism to have a natural network structures selection using this way of working:

- at the beginning each ant has a random network, random number of internal layers (for 1 to 2) having a random number of neurones (for now: for 5 to 50 for the first one and 5 to 30 for the second one if exist) and random synapses coef values. The input and output layers have all the same number of neurones in all ants, but it could be changed in future versions.
- on regular basis (1000 or 10000 ticks), for each ant:
  - compute the maximum distinct decisions the network is able to take (let's call it `maxDecision`)
  - compute good decision rate: [number of decision on the period which raise the happiness: `[reinforced decision number on the period] / [total number of decision on the period]`, let's call it `good decision rate`
- on regular basis (100000 ticks) during the simulation:
  - find the best network of the nest for each ant kind (worker, soldier) (the one having first the best `maxDecision` and if equal the one having the best `good decision rate`)
  - for each ants:
    - if its network `maxDecision` is lower for a factor 2 than the best network of its kind and if equal or lower factor 1, if its `good decision rate` is lower for 10% than the best network of its kind then: duplicate the best network of its kind with all it trained synapse coefficients and set it to a considered ant

this time, the network structures are not chosen anymore, the best ones will emerge naturally

Currently the version 0.0.1 (spread) is achieved and works with a good level of network convergence and about 80% of good decisions, the version 0.0.2 (foods and pheromones) works well, but it takes more time to reach the right percentage of good decision. version 0.0.3 have difficulties to converge due to the short life duration of the ants.
To counter-balance this effect, when a new ant is created (worker or soldier), 50% of the time, the nest give a copy of the best network of the ant kind and 50% of the time a random one (% can be adjust)
This way the best trained network is not lost if its associated ant died. It benefits to the others.


# Ant networks structures

On version 0.0.1, the input and output layer have 8 neurons each.
On version 0.0.2, the input layer has 24 neurons and the output 8.
On verison 0.0.4, the input layer has 32 neurons and the output 8.

In version 0.0.3, the networks have only one internal hidden layers, having for 5 to 50 neurons. Two internal layers networks never success to over take one internal layer network. Could be updated in future version if needed, but it speed up the first network selections stage (100000 to 200000 first ticks)

## input layer

On version 0.0.3:
- neurons 1 to 8 are dedicated to friend ants detection
- neurons 9 to 16 are dedicated to foods detection
- neurons 17 to 24 are dedicated to pheromones detection
- neurons 25 to 32 are dedicated to hostile ants detection


The number of neurons in a single detection slot, 8 neurons here, define the precision of the ant vision.

each neuron of a slot is associated to the ant vision quadrant:
The neuron 0 is associated to the quadrant 0 to 45 degrees
The neuron 1 is associated to the quadrant 45 to 90 degrees
The neuron 2 is associated to the quadrant 90 to 125 degrees
and so on...

The server is able to run with any precision number, 8 appeared to be enough and fast, but the server can be re-build to run with 12, 16 or 32 neurons per slot if needed.

so neurons:
- 0 is associated to the detection of friend ants in the quadrant 0 to 45 degrees
- 9 is associated to the detection of foods in the quadrant 0 to 45 degrees
- 17 is associated to the detection of pheromones in the quadrant 0 to 45 degrees
- 25 is associated to the detection of hostile ants in the quadrant 0 to 45 degrees
and
- 1 is associated to the detection of friend ants in the quadrant 45 to 90 degrees
- 10 is associated to the detection of foods in the quadrant 45 to 90 degrees
- 18 is associated to the detection of pheromones in the quadrant 45 to 90 degrees
- 26 is associated to the detection of hostile ants in the quadrant 45 to 90 degrees
and so on...

The first 8 neurons in the input layer activate itself when another ants appear in its quadrant, more the ant is close more the neuron is activated from 0 to 1.

The 9 to 15 neurons in the input layer activate itself when foods appear in its quadrant, more the food is close more the neuron is activated from 0 to 1.

The 16 to 24 neurons in the input layer activate itself when pheromones appear in its quadrant, more the food is close more the neuron is activated from 0 to 1.

The 25 to 32 neurons in the input layer activate itself when hostile ants appear in its quadrant, more the food is close more the neuron is activated from 0 to 1.

An ant has a maximum distance of vision, so an input neuron starts to activate itself when another ant or food enters inside the circle of vision and in the associated neuron quadrant.
the activated neuron value start to 0 when object detected is just at the edge of the vision, to 1 when the object is exactly at the same place then the ant


## output layer

The number of neurons in the output layer define the precision of the ant move. Each neuron drive a direction
neuron 0: move to 0 degree
neuron 1: move to 45 degree
neuron 2: move to 90 degree
...

The server is able to run with any precision number, 8 appeared to be enough and fast, but the server can be re-build to run with 12, 16 or 32 neurons in the output layer if needed.

When input layer is set and propagated through the network, then the output layer is used to decide to direction to take, see chapter `take decision on output`


# take decision on output

The way to take a direction decision using the network output is the following:

on regular basis, each ant compute a statistical distribution of the decisions taken using its network outputs.
This distribution is an array of int with for each possible decision (direction) the number of time the decisions has been taken during the period.
Let's call it `decisionSum` array. It's updated every 1000 ticks to always reflect the up to date situation.

The following algorithm is used to compute the decision using network output:
  - get the index of the neuron having the maximum value in the output layer, this index is the default decided direction
  - if the value of the decisionSum array at index-1 is lower than 100, then the taken decision become index-1
  if the value of the decisionSum array at index+1 is lower than 100, then the  taken decision become index+1

Given that index+1 and index-1 are close than the regular default decision, the result is not changed too much, but it gives a change to these decisions to be train and reinforced even if they are not the regular decisions, raising this way the discrimination capability of the network.

This counter-balance the flattering effect given by the bad decision fading. see chapter `Fade decision`


# Fade decision

If we only reinforce the good decisions, the ones which raise the ant happiness, it takes too much time to the networks to converge and some of them can't converge because their random initial synapse values lead to only bad decisions.

It's necessary to use also the bad decisions to train network. We have plenty of them more than good ones, especially at the beginning of the simulation.

For instance, we get an entry with input value [0, 0.5, 0, 0, 0, 0, 0, 0, 0] and the direction 2 decision is taken after having propagated it,
The loop after, this decision appears to be bad, because after the move the ant happiness lower, so we need to fade this decision training the network with a computed fade output for the decision (direction).

The output corresponding to the decision 2 is [0, 0, 1, 0, 0, 0, 0, 0, 0], so we can train the network using retro-propagation algorithm with the couple:

input:  [0, 0.5, 0, 0, 0, 0, 0, 0, 0]
output: [x, x, 0, x, x, x, x, x, x]

where x is a parameter, for now set to 0.3, but it has a great influence on the network capability to converge and the best should be to have it random and let the best values emerge by themselves using the network purge process (on some next versions)

Training several times the network this way with bad decisions, the network is forced to "forget" little by little them.


# The happiness of a ant

To train a network, it needs samples, couples of input-output known are right for the purpose of the network.

They aren't such samples in the ant networks case. In fact it could have, it's possible to compute them especially outside the simulation in order to train networks on very effective way, but it's not the purpose of this project.

In this project, it's supposed that compute artificial samples is not possible. It's in this case that the neuron networks are interesting, when the good answers emerge, not because of the validity of samples created by an external algorithm, but because networks are driven using a high level parameter, ants happiness here which is relevant.

To achieve that:

First, we get input data in the simulation context itself, the ant visions, so we have plenty inputs at out disposal just letting ants move.

Second, after network propagation of this inputs, we need to assess if the output is a good or a bad one in order to reinforce or fade them.

To do that, for each ant, at each ant move, the happiness of the ant is computed. if after a move the happiness is lower than the previous one, the output is considered as bad, if the happiness is higher the output is considered as good. if it's the same, decision is not considered at all.

Then the network convergence become driven by happiness, a unique high level parameter we can define as we need.

The happiness should be a function of the network entry layer neuron values only, if an external parameter is used in the computation and if this parameter can change independently than the input values, the network won't be able to converge. There is no magic, the network needs to correlate entries values only with outputs values to be trained a way making it able to converge.

In version 0.0.1:
The happiness of a ant is as greater as it sees no other ant. The happiness will lower more and more, when other ants enter in its vision field and they are close.

The computation of the happiness for each ant is:
- considering the ants detected by the input layer, sum the power 2 of their distances

In the version 0.0.2:
The computation of the happiness for each ant is:
- considering the food detected in the input layer, sum the power 2 of their distances
- if this sum if > 0 stop there, happiness is this sum
- considering the pheromones detected in the input layer, get the power 2 of the distance of the one having the lower level (pheromone evaporate themselves with time and are less detected with time)
- if this result if > 0 stop there, happiness is this value
- considering the friend ants detected in the input layer, sum the power 2 of their distances
- happiness is this sum

In the version 0.0.3:
it's needed to consider the two ant kinds:
- workers happiness computing:
  - considering the hostile ants detected in the input layer, sum the power 2 of their distances
  - if this sum if > 0 stop there, happiness is this sum
  - considering the food detected in the input layer, sum the power 2 of their distances
  - if this sum if > 0 stop there, happiness is this sum
  - considering the pheromones detected in the input layer, get the power 2 of the distance of the one having the lower level (pheromone evaporate themselves with time and are less detected with time)
  - if this result if > 0 stop there, happiness is this value
  - considering the friend ants detected in the input layer, sum the power 2 of their distances
  - happiness is this sum
-soldier happiness computing:
- considering the hostile ants detected in the input layer, sum the power 2 of their distances
- if this sum if > 0 stop there, happiness is this sum
- considering the friend ants detected in the input layer, sum the power 2 of their distances
- happiness is this sum

At the end soldier could have only 16 neurons in their entry layer, concidering the foods and pheromones ones are not used. to be decided...


# Install and build

This project uses a server written in Go and an Angular 4 client.

Pre-requisites are:
- golang 1.8 installed (with golint, to install golint: `go get -u github.com/golang/lint/golint`)
- git installed
- make

To install and build:

- clone the git project: https://github.com/freignat91/mlearning on your $GOPATH and execute these commands:
  - cd $GOPATH/src/github.com/freignat91 (create it)
  - git clone git://github.com/freignat91/mlearning
- build the project, executing this commands:
  - cd $GOPATH/mlearning
  - make

Then there are two executables in your $GOPATH/bin directory:
  - ml: the mlearning command line  
  - mlearning: the server



# usage

## server

To start server execute the command: `mlserver`
($GOPATH/bin should be in your $PATH, if not execute $GOPATH/bin/mlserver)

## UI

To see the UI, open a Chrome (tested only on Chrome for now) and enter url: localhost:3001

Then you can:
  - start/stop the simulation and use "next step" button to move tick after tick and see result
  - speed up/down the server (down to be able to see the moves, up to let train the network faster)
  - Button "clear group", remove all existing food groups
  - two click modes are possibles:
    - click on "Select ant" button and then on the graphic space to select an ant
    - click on "AddFoods" button and then on the graphic space to add a food group
  - menu "File":
    - "Restart": (re)start a new simulation from the beginning
    - "Export sample": export max 10000 trained sample of the selected ant to file (the file is created on server side ./test/testant.json)
  - checkboxes:
    - contact circles: if checked, show circles when ants of the same nest are in contact (less than vision length / 4) and when there in scoot mode (not car
      rying food, not panic, not fighting, not trace up pheromone)
    - fight circles:if checked, show circles when ants are fighting
    - Food renew:
      - not checked: stop the food to be replaced in the existing groups when they reach the nest
      - checked: (re)start to replace food in the existing groups when they reach the nest
    - don't panic: if checked, make workers don't panic, so no more soldiers can be created. Clicked at the beginning, nests feed harmoniously, showing that soldiers are not useful when ants don't panic.

back or red circles, depending on the ant color show the vision circle of the selected ant.
The green spots are the foods and the back/grey dash are the pheromones.


The UI shows the graphical simulation on the left and information on the right

where:
- from the beginning, total and for the selected ant:
  - Timer is the tick number (one tick compute one move for all the ants) and the number of ticks per second
  - update networks: the updated (copied) networks
  - train number: the number of time the networks have been trained
  - decision rate: the good decisions % (good decisions / total decisions)
- for the current period (2 sec), total and for the selected ant:
  - contact: the number of ant contacts (less than the max vision length / 2)
  - train number: the number of time the networks have been trained
  - decision: the number of decisions taken
  - positive decision reinforcement: the number of trains after a good decision
  - negative decision fading: the number of fading train after a bas decision
  - update networks: the updated (copied) networks
  - train number: the number of time the networks have been trained
  - decision rate: the good answers % (good decisions / total decisions)
- for networks assessment, worse or best:
  - ant id of network (clicking on it, it selects it)
  - structure: structure of the networks (number of neurons per layer)
  - distinct decisions: distinct number of decision the network can taken
  - decision rate: the good decisions % of the network (good decisions / total decisions)
- for nests back and red:
  - ressources: ability to create new ant (4 point for a soldier, 1 for a worker)
  - life: sum of all the life of the nest ants
  - worker: number of workers
  - soldier: number of soldier


## command line

There are command line to:
- work with a neuron network, create it, train it, verify result of it, independently than the simulation
- work with a network of a given ant in the simulation, select it, test it, train it, ...

Execute the commande: `ml --help` or `ml network --help` to get help

Commandes list:

### ml network create x1 x2 x3 ... xn

create a new neuron network having:
  - x1 neurons in the input layer
  - xn neurons in the output layer
  - x1 to x(n-1) neurons in the hidden layers

this network become the current one

### ml network display

display information on the current network, layers, neuron values and if the option `--coef` is set, the coef values

### ml network propagate val1, val2 .. valn

push the values val1 to valn in the neurons of the input layer of the current network and propagate them.
display the outpit values

### ml network backPropagate val1, val2, ....

using last propagate values, push values val1 to valn in the output layer of the current network and train it, to adjust synapses coef using back propagation algorithm

### ml network loadTrainFile path

load a training json file having couples of input-output.
See ./tests directory for examples

For instance:
```
{
  "name": "testxor",
  "layers": [2, 3, 1],
  "data": [
    {
      "in": [1, 1],
      "out": [0]
    },
    {
      "in": [1, 0],
      "out": [1]
    },
    {
      "in": [0, 1],
      "out": [1]
    },
    {
      "in": [0, 0],
      "out": [0]
    }
  ]
}
```

### ml network train name -n x -c

train the current network using the loaded sample data (previously loaded) named `name`.
execute all the samples, `n` time
create the network first, if option `-c` is set. the network is created even without this option if no current network exist

### ml network select nestId antId

get the network of the ant `antId` in the nest `nestId` and set it as current network
if nestId = "best", set the best network as current
if nestId = "worse" set the worse network as current
if not argument, set the current ant (selected in UI) as current

### ml network test nestId nestId antId

set the network of the ant `antId` in the nest `nestId` as current and test it
if nestId = "best", set the best network as current and test it
if nestId = "worse" set the worse network as current and test it
if not argument, set the current ant (selected in UI) as current and test it

test it, means propagate the main possible entry to see how the network converge and is able to distinct decision:

```
ml network test best
Test network: [8 22 8]
[ 1 0 0 0 0 0 0 0 0 ] => [ 0.05 0.03 0.03 0.03 (0.05) 0.84 0.04 0.03 ] max=5/4 diffMax=0.70
[ 0 1 0 0 0 0 0 0 0 ] => [ 0.03 0.02 0.03 0.03 0.04 (0.03) 0.83 0.03 ] max=6/5 diffMax=0.70
[ 0 0 1 0 0 0 0 0 0 ] => [ 0.04 0.02 0.05 0.03 0.03 0.04 (0.04) 0.78 ] max=7/6 diffMax=0.66
[ 0 0 0 1 0 0 0 0 0 ] => [ 0.89 0.01 0.03 0.01 0.03 0.02 0.01 (0.02) ] max=0/7 diffMax=0.76
[ 0 0 0 0 1 0 0 0 0 ] => [ (0.09) 0.10 0.65 0.09 0.07 0.08 0.08 0.08 ] max=2/0 diffMax=0.50
[ 0 0 0 0 0 1 0 0 0 ] => [ 0.02 (0.81) 0.06 0.02 0.03 0.03 0.02 0.02 ] max=1/1 diffMax=0.68
[ 0 0 0 0 0 0 1 0 0 ] => [ 0.02 0.01 (0.03) 0.84 0.03 0.01 0.02 0.01 ] max=3/2 diffMax=0.71
[ 0 0 0 0 0 0 0 1 0 ] => [ 0.03 0.03 0.02 (0.03) 0.87 0.03 0.03 0.02 ] max=4/3 diffMax=0.74
[ 0 0 0 0 0 0 0 0 1 ] => [ 0.07 0.05 0.16 0.12 (0.09) 0.06 0.05 0.06 ] max=2/4 diffMax=0.08
Match rate:0.00000 tot:5.52641 distinct=8
```

where:
- max=x1/x2, x1 is the index of the max value, x2 the theoretical best answer
- diffMax: the difference between the max and the average other value
- tot: the su; of all deffMax
- distinct: the distinct direction the network is able to output



## License

mllearning is licensed under the Apache License, Version 2.0. See https://github.com/freignat91/mlearning/blob/master/LICENSE
for the full license text.
