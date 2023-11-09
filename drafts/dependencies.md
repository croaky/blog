# Dependencies

WIP: follow-up to "Deep modules, simple interfaces"

In 2023, it is unlikely that many engineering teams will build at a lower level
of abstraction than Twilio for telephony, Stripe for moving money, or Cloudflare
for DDoS protection.

However, we are in the early days of AI. Is it less clear that an API is too
high level of an abstraction for handling tasks that use a large language model?

An engineer working in a module must understand
both the interface and implementation of that module,
plus the interfaces of modules invoked by that module.

Modules work together by calling each others' functions or methods.
There will be dependencies between modules.

To manage dependencies, we think of each module in terms of
its interface and its implementation.
The interface consists of what the caller must know,
describing what the module does but not how.
The implementation consists of the code that
carries out the promises made by the interface.

The smaller and simpler the interface,
the less complexity it imposes on the rest of the system.
Complexity is anything related to the structure of the system
that makes it hard to understand and modify.

Complexity is what a developer experiences. Symptoms are:

- Change amplification: A simple change requires many pieces of code to change.
- Cognitive load: How much a developer must know to complete a task.
- Unknown unknowns: Not obvious what must change to complete a task.

Change amplification is tedious.
Cognitive load increases the cost of change.
Unknown unknowns are the worst:
there is something you need to know,
but you don't know how to find out about it,
and may not even know if there is an issue.

## Abstractions

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
and to minimize the amount of information that is important.

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

Information hiding reduces complexity by:

- Simplifying the interface to a module
- Making it easier to evolve the system

When designing a module,
if you can hide more information,
you should be able to simplify the module's interface,
and this makes the module deeper.

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

## General purpose modules are deeper

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
