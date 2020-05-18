# Routes

In both server-side APIs and front-end web or mobile apps,
we need to map "routes" to functions.

The simplest and most maintainable solution to this problem
was first established (to my knowledge) by
[Sinatra](http://sinatrarb.com/intro.html)
and later picked up by frameworks like
[Express](https://expressjs.com/en/guide/routing.html)
and conventions in Go codebases.

## Server-side APIs

> In Sinatra, a route is an HTTP method paired with a URL-matching pattern.
> Each route is associated with a block:

```rb
get "/" do
  # ...
end

post "/" do
  # ...
end
```

In my opinion, this is the key attribute of a maintainable routing system:

> Routes are matched in the order they are defined.
> The first route that matches is invoked.

Secondarily, it's nice to have a conventional `routes` file of some kind
that defines these routes.

## Front-end web

Routes in front-end frameworks like React that I've worked with
have tended to not have a clear routing structure.

## Front-end mobile
