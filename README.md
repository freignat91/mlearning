# Ants Machine learning

mlearning v0.0.1

# Purpose

Use multi-layer neuron network to drive ant behaviours, each ant having its own network.

Display graphical representation of the ant behaviours and information to monitor how they evolve and how their network converge.

Use the virtual ant context to train the networks, train the networks permanently using what the ants see and the consequences of theirs actions.

version 0.0.1: Ants should be able to spread the virtual space in order to cover all the available space witout having contact between them. The ants avoid each other

version 0.0.2: Foods will appear in the space, ants should be able to get them and bring back them to the nest, letting a pheromone path

version 0.0.3: Ants should be able to trace back the a pheromone path to find food sources

version 0.0.4: with two nests, ants should be able to fight agains ants of the other nest.


# Way of working

Ant network are train using regular gradient retro-propagation algorithm. The point is to be able to define what is the right output of a network for a given entry.

For that let's consider the consequence of an ant decision, using the following algorithm:

- for each simulation tick and for each ant:
  - build the network entry using what the ant see (see chapter `ant network structure`)
  - propagate the entry in the network
  - considering the network output, take a decision about the direction the ant should move (see chapter `take decision on output`)
  - move the ant to the chosen direction
  - compute the happiness of the ant (see chapter `the happiness of a ant`)
  - if the happiness is the same than the one of previous loop, do nothing and keep the direction as it is.
  - if this happiness is greater than the one of previous loop, then consider that the decision is good and train the network using the previous entry and the output corresponding the chosen direction.

This loop appeared not to be enough in order to make the network converging to stable results. It needed to not only reinforce the good decisions, but to fade the bad ones.

Then it's possible to add the following step:

- if the computed happiness is lower than the previous one, train the network with the previous entry and a fading output computing using the chosen direction (see chapter `fade decision`)

With this additional rule, the network started to converge well, but a new issue appeared:

Because of the fading, the networks result become more and more poor, meaning, the number of the distinct possible output decisions lower. For instance a network shows only 3 directions no matter the entries, where it should be able to show 8.

At the end, all the ants finished to move to only 2 directions, even if they reached the assigned goal: they cover all the virtual space and the contacts number lower to 0, it's not good for the next version purpose.

Hopefully, it's possible to enrich the network decision capability to counter-balance the fade effect. To do that, let's use the following algorithm:

On regular basis (every 1000 or 10000 ticks), for each ant:
- compute the statistical network output distribution, for each possible decision: the number of time the decision has been taken in the period
- when decision has to be taken using network output:
  - compute the regular decision as usual (see chapter `take decision on output`)
  - look at the immediate other direction (+1 and -1 regarding the regular chosen one) and if they haven't been take more than a given time (100 for now), then take this decision instead of the regular one.
  - This decision will be as usual reinforced or fade considering its consequence

This way, we train networks on decision it doesn't show too much on its regular way. If the "forced" decisions appear to be good, they will be naturally reinforced and then their statistical distribution will be better (>100) and then they become regular decisions.

Then, the networks are now able to converge to use about all their decision capability.

Now, there is another issue: We don't all which network structure it should be used. 3 layers, 4 layers? how many neurones by layer?

The fist tested network were empirically set at the beginning finding that a network 8-7-8 is not bad for the purpose (see chapter `ant network structure`), but perhaps 8-30-100-8 is better...

Then, it's possible to add a mechanism to have a natural network structure selection using this way of working:

- at the beginning each ant has a random network, random number of internal layers (for 1 to 2) having a random number of neurones (for now: for 5 to 50 for the first one and 5 to 30 for the second one if exist) and random synapses values. For now, the input and output layers have the same number of neurones for all ants, but it could be changed for future versions.
- on regular basis (1000 or 10000 ticks), for each ant:
  - compute the maximum distinct decisions the network is able to take (let's call it `maxDecision`)
  - compute good decision rate: [number of decision on the period which raise the happiness] / [total number of decision on the period], let's call it `good decision rate`
- on regular basis (100000 ticks) during the simulation:
  - find the best network of the nest (the one having first the best `maxDecision` and if equal the one having the best `good decision rate`)
  - for each ants:
    - if its network `maxDecision` is lower for a factor 2 than the best network one and if not, if its `good decision rate` is lower for 20% than the best network one then: duplicate the best network with all it trained synapse coefficients and set it to a considered ant

this time, the network structures are not chosen anymore, the best ones will emerge naturally

All this way of working works for the initial ant spread purpose, but should also works for all the others:  find food, trace up pheromones, fight, ...

Currently the version 0.0.1 (spread) is achieved and works enough well to move to version 0.0.2, but not sure to keep one network per ant, perhaps one network per task with a meta-network to activate the task networks considering a meta happiness. to be see...


# Ant networks structures

On version 0.0.1, the input and output layer have 8 neurons each. The networks have one or two internal hidden layers, having for 5 to 50 neurons for the first one and 5 to 30 for the second one. These numbers was arbitrary chosen at the beginning and proven not so bad by tests.

## input layer

The number of neurons in the input layer define the precision of the ant vision.
each neuron is associated to the ant vision quadrant.
The neuron 0 is associated to the quadrant 0 to 45 degree
The neuron 1 is associated to the quadrant 45 to 90 degree
The neuron 2 is associated to the quadrant 90 to 125 degree
and so on...

The server is able to run with any precision number, 8 appeared to be enough and fast, but the server can be re-build to run with 12, 16 or 32 neurons in the input layer if needed.

A neuron in the input layer activate itself when another ant appears in its quadrant, more the ant is close more the neuron is activated from 0 to 1.

An ant has a maximum distance of vision, so an input neuron starts to activate itself when another ant enters inside the circle of vision and in the associated neuron quadrant

It'll be the same for food, pheromone, hostile ant, for next version, but each kind of object to detect will have its own circle of neurons. So 8 for friend ants, 8 for food, 8 for pheromones, 8 for hostile ants detection.

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

on regular basis (every 1000 ticks), each ant compute a statistical distribution of the decision made using their network output.
This distribution is an array of int with for each possible decision (direction) the number of time the decisions has been take during the period.
Let's call it decisionSum array. It's updated every 1000 ticks to always reflect an up to date situation

The following algorithm is used to compute the decision using network output:
  - get the index of the neuron having the maximum value in the output layer, this index is the default decided direction
  - if the value of the decisionSum array at index-1 is lower than 100, then the taken decision become index+1
  if the value of the decisionSum array at index+1 is lower than 100, then the  taken decision become index-1

Given that index+1 and index-1 are pretty close than the regular default decision, the result is not changed too much, but it gives a change to a decision which is not the regular network one to be train and reinforced raising this way the discrimination capability of the network.

This counter-balance the flattering effect given by the bad decision fading. see chapter `Fade decision`


# Fade decision

If we only reinforce the good decisions, the ones which raise the ant happiness, it takes too much time to the networks to converge and some of them can't converge because their random initial synapse values lead to only bad decisions.

The idea, is to use also the bad decisions to train network. We have plenty of them more than good ones, especially at the beginning of the simulation.

For instance, we get an entry with input value [0, 1, 0, 0, 0, 0, 0, 0, 0] and the direction 2 decision is taken after having propagated it,
The loop after, this decision appears to be bad, because after the move the ant happiness lower, so we need to fade this decision training the network with the opposite of what we got as output layer the first time.

The ideal output corresponding to the decision 2 is [0, 0, 1, 0, 0, 0, 0, 0, 0], so we can train the network using retro-propagation algorithm with the couple:

input:  [0, 1, 0, 0, 0, 0, 0, 0, 0]
output: [x, x, 0, x, x, x, x, x, x]

where x is a parameter, for now set to 0.3, but it has a great influence on the network capability to converge and the best should be to have it random and let the best values emerge by themselves using the network purge process (on some next versions)

Training several times the network this way, the network is forced to "forget" the decision.



# The happiness of a ant

To train a network, it needs samples, couples of input-output known are right for the purpose of the network.

They aren't such samples in the ant networks case. In fact it could have, it's possible to made them especially in order to train networks on very effective way, but it's not the purpose of this project.

In this project, it's supposed that, made artificial samples is not possible. It's in this case that the neuron networks are interesting, when the good answers emerge, not because of the validity of samples created by an external algorithm, but because we drive them using a high level parameter, ants happiness here.

To achieve that:

First, we get input data in the simulation context itself, in version 0.0.1 input of the networks are the ant visions so we have plenty inputs at out disposal just letting ants move.

Second, after network propagation of this inputs, we need to assess if the output is a good one in order to reinforce or fade network outputs.

To do that, for each ant, at each ant move, the happiness of the ant is computed. if after a move the happiness is lower than the previous one, the output is considered as bad, if the happiness is higher the output is considered as good.

Then the network convergence become driven by this unique high level parameter we can define as we need.

The happiness should be a function of the network entry layer neuron values only, if an external parameter is used in the computation and if this parameter can change independently than the input values, the network won't be able to converge. There is no magic, the network needs to correlate entries only values with outputs to be train a way making it able to converge.

For the 0.0.1 version, the happiness of a ant is as greater as it sees no other ant. The happiness will lower more and more, when other ants enter in its vision field and they are close.

On the version 0.0.1 the computation of the happiness for each ant is:
- find all the ants inside the field vision
- sum the distance of these ants (in fact the power 2 of the distance, no need to compute the sqrt)

For the 0.0.2 the happiness will be compute the same way but inverse for food, more the food is close and more there are foods, more the happiness of a ant raise.


# Install and build

This project use a server written in Go and an Angular 4 client.

Pre-requisite are:
- golang 1.8 installed
- git installed

To install and build:

- clone the git project: https://github.com/freignat91/mlearning on your $GOPATH and execute these commandes:
  - cd $GOPATH
  - git clone git@github.com:freignat91/mlearning
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
  - click on an ant on the graphical representation select it
  - export max 10000 trained sample of the selected ant to file (the file is created on server side ./test/testant.json)

The UI shows the graphical simulation on the left and information on the right:

Information are for instance:
```
Timer	                              2030000	   30880 t/s
From the beginning                  Global     Selected ant
Updated network	                    1578	     13
Train number	                      2819070	   27606
gRate	                              44.85%	   47.43%

Average on period
Contact	                            0	         0
Train number	                      345	       370
Decision	                          437	       510
Positive decision reinforcement	    345	       370
Negative decision fading	          91	       140
Updated network	                    95	       0
gRate	                              79.03%     72.55%

Networks                            Worse      Best
Id	                                76	       17
Structure	                          [8 22 8]	 [8 22 8]
DirCount	                          4 	       8
gRate	                              75.51%	   81.55%
```

where:
- from the beginning, total and for the selected ant:
  - Timer is the tick number (one tick compute one move for all the ants) and the number of ticks per second
  - update networks: the updated (copied) networks
  - train number: the number of time the networks have been trained
  - good decision rate: the good decisions % (good decisions / total decisions)
- for the current period (2 sec), total and for the selected ant:
  - contact: the number of ant contacts (less than the max vision length / 2)
  - train number: the number of time the networks have been trained
  - decision: the number of decisions taken
  - positive decision reinforcement: the number of trains after a good decision
  - negative decision fading: the number of fading train after a bas decision
  - update networks: the updated (copied) networks
  - train number: the number of time the networks have been trained
  - good decision rate: the good answers % (good decisions / total decisions)
- for networks assessment, worse or best:
  - ant id of network (clicking on it, it selects it)
  - structure: structure of the networks (number of neurons per layer)
  - distinct decisions: distinct number of decision the network can taken
  - good decision rate: the good decisions % of the network (good decisions / total decisions)


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
