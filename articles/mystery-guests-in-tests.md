# Mystery Guest

> You're having trouble understanding the behavior a test is verifying.

This testing anti-pattern is called [Obscure Test].
A common cause of Obscure Test is [Mystery Guest]:

[obscure test]: http://xunitpatterns.com/Obscure%20Test.html
[mystery guest]: http://xunitpatterns.com/Obscure%20Test.html#Mystery%20Guest

> The test reader is not able to see the cause and effect between fixture and
> verification logic because part of it is done outside the Test Method.

The impact is:

- Tests don't fulfill the role of [Tests as Documentation].
- You may have [Erratic Tests] which result don't pass during every test run or
  pass in the test environment but not in production.

[tests as documentation]: http://xunitpatterns.com/Goals%20of%20Test%20Automation.html#Tests%20as%20Documentation
[erratic tests]: http://xunitpatterns.com/Erratic%20Test.html

## Identify the Mystery Guest

In the test suite of web applications,
the Mystery Guest might be one of these smells:

- a fixture function
- a fixture file such as XML or JSON saved from an HTTP response
- a variable defined at the top of a long code block
- a variable defined up multiple levels of nesting
- a variable without an [Intention-Revealing Name]

[intention-revealing name]: http://c2.com/cgi/wiki?IntentionRevealingNames

Some smells have direct solutions:

- shorten long blocks
- prefer flat over nested tests
- rename variables

Other smells require more thought.

## Replace shared fixtures with fresh fixtures

Are your tests using shared fixtures like this?

```yaml
dan:
  name: Dan
  role: developer
  location: San Francisco

phil:
  name: Phil
  role: designer
  location: Boston
```

Is it important in a specific test that the fixtures used
are role-based (developer) or location-based (San Francisco)?
If the former is true,
we should rename the fixtures in order to reveal the intention:
`users(:developer)` and `users(:designer)`.

What we call "creating objects with factories",
others call creating [Fresh Fixtures]
built using [Inline Setup].

[fresh fixtures]: http://xunitpatterns.com/Fresh%20Fixture.html
[inline setup]: http://xunitpatterns.com/Inline%20Setup.html

> Each test method creates a test fixture for its own private use.

This meets our goal for the test reader to better see the cause and effect
between fixture and verification logic.

## Set up only relevant information

[Irrelevant Information] is another cause of Obscure Test.
Again, Factory Girl is helpful.
Consider this setup:

[irrelevant information]: http://xunitpatterns.com/Obscure%20Test.html#Irrelevant%20Information

```ruby
context "user account exists with a matching Facebook uid" do
  setup do
    @uid  = 1234567
    @user = create(:user, fb_user_id: @uid)
  end
  ...
end
```

We explicitly specify only the attribute on user that matters for this test.
We use Factory Girl to create only a user with valid data,
and we name a `@uid` variable to express intent for the otherwise
[Magic Number] `1234567`.

[magic number]: http://en.wikipedia.org/wiki/Magic_number_(programming

## External resources

Sometimes we need to assert the contents of a file are what we expect.
For example, say we've generated the file by dumping some XML from a web service
which we'll then stub out during test runs, using the file as a proxy.

```ruby
context "the job XML from the web service" do
  should "include the recruiter email" do
    @xml_job = IO.read("test/fixtures/jobs/1.xml")
    recruiter_xml = "<recruiter>recruiter@example.com</recruiter>"
    assert_contains @xml_job, recruiter_xml
  end
end
```

In this case, `test/fixtures/jobs/1.xml` is called a [Prebuilt Fixture].

[prebuilt fixture]: http://xunitpatterns.com/Prebuilt%20Fixture.html

To truly avoid a Mystery Guest here, the XML could be inline.

## Conclusion

Mystery Guest is memory trick to think about how code will be
read and understood by humans.
Can the intended behavior of a subset of the system be understood at a glance?
Or is there a Mystery Guest clouding our understanding?
