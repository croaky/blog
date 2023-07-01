# Philosophy Of Software Design

One of my favorite books on software is
[A Philosophy of Software Design](https://amzn.to/2OQkBEQ).
It provides useful definitions of software complexity and simplicity,
how you can tell if a system is unnecessarily complex,
what causes a system to become complex,
and how to design software systems to minimize their complexity.

I especially love its sections about "deep" modules.

## Definition

Complexity is anything related to the structure of a software system
that makes it hard to understand and modify the system.

It might be hard to understand how a piece of code works.
It might take a lot of effort make a small change.
It might not be clear which pieces of the system to modify to make the change.
It might be difficult to fix one bug without introducing another.

Complexity is what a developer experiences.
Complexity is more apparent to readers than writers.

## Symptoms

How to notice complexity:

- Change amplification: A simple change requires many pieces of code to change.
- Cognitive load: How much a developer must know to complete a task.
- Unknown unknowns: Not obvious which code must change to complete task.

Change amplification is tedious.
Cognitive load increases the cost of change.
Unknown unknowns are the worst of the types:
there is something you need to know,
but you don't know how to find out about it,
and may not even know if there is an issue.

Non-symptom: number of lines of code.

## Causes

Complexity is caused by:

- Dependencies
- Obscurity

A dependency exists when a piece of code
cannot be understood and changed in isolation:
it relates in some way to other code,
and the other code must be considered (and maybe modified)
if the given code is changed.
Dependencies are fundamental to software and cannot be eliminated.
We often intentionally introduce dependencies.
Some goals of software design are:

- Reduce dependencies
- Make dependencies simple and obvious

Obscurity occurs when important information is not obvious.
Unclear variable/method/parameter names might require
looking up documentation or implementation of dependencies.
Inconsistency is a major contributor to obscurity.
Obscurity can be an issue of inadequate documentation
but it is also a design issue.
If a system has a clean and obvious design,
it will need less documentation.
The best way to reduce obscurity is to simplify the system design.

## Modules should be deep

Software can be decomposed into a collection of modules
that are relatively independent.
Modules can be classes, subsystems, or services.

Ideally, an engineer could work on one module without worrying about another.
Modules, though, work together by calling each others' functions or methods.
There will be dependencies between modules.

To manage dependencies, we think of each module in terms of
its interface and its implementation.
The interface consists of what the caller must know,
describing what the module does but not how.
The implementation consists of the code that
carries out the promises made by the interface.

An engineer working in a module must understand
both the interface and implementation of that module,
plus the interfaces of modules invoked by that module.

The best modules are those whose
interfaces are much simpler than their implementations.

## What's in an interface?

The interface to a module contains formal and informal information.

The formal parts are specified explicitly in the code,
and some of these can be checked for correctness by the programming language.
The formal interface of a class consists of the signatures for
all of its public methods, plus names and types of its public variables.

The informal parts are not specified in a way that can be enforced
by the programming language or other tooling.
For example, a function might delete a file named by one of its arguments.
This information needs to be understood by the engineer,
and is therefore part of its interface.
It can only be described by its documentation.

## Abstractions

A microwave oven contains complex electronics
converting alternating current into microwave radiation and
distributing that radiation throughout the cooking cavity,
but abstracts an interface of a few buttons to control timing and intensity.

An abstraction is a simplified view of an entity,
which omits unimportant details.

A detail can only be omitted from an abstraction if it is unimportant.
An abstraction can go wrong in two ways:

1. It can include details that are not important.
   This makes the abstraction more complicated than necessary,
   increasing cognitive load on engineers using it.
2. It can omit details that are important.
   This results in obscurity.
   Engineers looking only at the abstraction
   will not have all the information they need to use it correctly.
   This is a "false abstraction": it appears simple but isn't.

The key to designing abstractions is understand what is important,
and to minimize the amount of information is important.

## Deep modules

"Deep" modules are those that provide powerful functionality
yet have simple interfaces.

Consider `===` lines the interface (cost: less is better)
and `|` lines the height (benefit: more is better):

```
===========
|         |
|         |
|         |
|         |   =============================
|         |   |                           |
|         |   |                           |
-----------   -----------------------------
Deep module   Shallow module
```

The benefit provided by a module is its functionality (depth).
The cost of a module is its interface (breadth).

A module's interface represents the complexity that
the module imposes on the rest of the system:
the smaller and simpler the interface,
the less complexity that it introduces.

An example of a deep module is the garbage collector in a language such as Go.
This module has no interface at all;
it works invisibly behind the scenes to reclaim unused memory.
Adding garbage collection to a system actually shrinks its overall interface,
since it eliminates the interface for freeing objects.

## Shallow modules

A shallow module is one whose interface is complex
in comparison to the functionality that it provides.

An extreme example:

```java
private void addNullValueForAttribute(String attribute) {
  data.put(attribute, null)
}
```

For managing complexity, this method makes things worse, not better.

- It offers no abstraction, all its functionality is visible through its interface.
- It is no simpler to think about the interface than the full implementation.
- It is more keystrokes than manipulating the `data` variable directly.
- It adds complexity (a new interface for engineers to learn)
  but provides no compensating benefit.

## Classitis and small methods

Conventional wisdom in programming is that classes should be small, not deep.
Programmers are told to break up larger classes into smaller ones.
The same is said for methods:

> methods longer than N lines should be divided into multiple methods

Classitis results in classes that are individually simple
but produce complexity from the accumulated interfaces.

## Information hiding

The most important technique for creating deep modules is information hiding.
Each module encapsulates knowledge of design decisions,
embedded in the implementation but not visible in the interface.

For example:

- How to store and access data
- What kind of data structure to use
- We'll use a particular JSON parser
- Pagination size
- Most files will be small

Information hiding reduces complexity by:

- Simplifying the interface to a module
- Making it easier to evolve the system

When designing a module,
if you can hide more information,
you should be able to simplify the module's interface,
and this makes the module deeper.

Hiding variables and methods in a class with `private`
isn't the same thing as information hiding.
Private elements can help hide information
but information about them can still be exposed
through public methods such as setters and getters.

## Information leakage

Information leakage occurs when a design decision
is reflected in multiple modules,
creating a dependency between them:
any change to that design decision
will require changes to all involved modules.

For example,
if two modules depend on a file format
so that one can write and the other can read,
they both depend on the file format.
Dependencies that aren't obvious through the interface are especially harmful.

If you encounter leakage between modules,
ask "How can I reorganize this code
so knowledge of this design decision
only affects one module?"

It might make sense to merge two modules.

Or, it might make sense to extract a new module
responsible only for that design decision.
This will only be effective if you can design a simple interface
that abstracts away the details.
If the new module exposes most of the knowledge through its interface,
then it won't provide much value:
it replaces back-door leakage with leakage through an interface.

## General purpose modules are deeper

Are you designing a general purpose or special purpose module?

General purpose means addressing a broad range of problems,
not only the ones that are important today.
Spend a bit more time up front to save time later.
It's hard to predict the future needs of software,
so general purpose might include facilities that are never needed.
If something is too general purpose,
it might not do a good job of solving the problem you have today.

So, maybe it's better to focus on today's needs,
building only what we know we need in a specialized way?
If we discover additional uses later,
we can refactor to make it general purpose.
This feels incremental.

The sweet spot is when a module's functionality reflects current needs
but its interface is general enough to support multiple uses.

## Different layer, different abstraction

Software is composed in layers.
Higher layers use facilities provided by lower layers.
Each layer should provide a different abstraction from layers above and below.

## Pull complexity downwards

Unavoidable complexity can be handled
either by users of a module or by internals of the module.

Since most modules have more users than developers,
it is more important for a module to have a simple interface
than a simple implementation.

A temptations of a module developer is to "let the caller handle it":

- raise an exception
- define a configuration parameter to control a policy

These approaches amplify complexity.
Many people must deal with a problem instead of one.
If a class throws an exception, every caller must handle it.
If a configuration parameter is exported,
every user may need to learn to set it.

## Better together or better apart?

The goal is to minimize overall system complexity.

To achieve this goal, we could divide the system
into a large number of small components.
The smaller the component, the simpler each component will be.

The act of dividing creates additional complexity that was not there:

- More components means more difficulty to distinguish between them
  and find the right one for each job.
- Additional code to manage.
- Separation: the code is farther apart than before.
  If the components are truly independent,
  this is good (the developer can focus on a single component at a time).
  If there are dependencies between them,
  this is bad (the developer needs to flip between them, or is unaware of them).
- Duplication.

Indications two pieces of code may be related:

- They share information. For example, they depend on syntax of a document.
- They overlap conceptually. For example: searching for a substring and
  case conversion are both string manipulation. Flow control and
  reliable delivery are both network communication.
- It is hard to understand one of the pieces without looking at another.
