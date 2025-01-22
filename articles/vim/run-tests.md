# vim / run tests

Test-driven development thrives on a tight feedback loop
but switching from editor to shell
to manually run specs is inefficient.

The [vim-test](https://github.com/vim-test/vim-test) plugin
exposes commands such as `:TestNearest`, `:TestFile`, and `:TestLast`,
which I bind to `<Leader>s`, `<Leader>t`, and `<Leader>l`.

Cursor over any line within an RSpec spec like this:

```ruby
describe RecipientInterceptor do
  it 'overrides to/cc/bcc fields' do
    Mail.register_interceptor RecipientInterceptor.new(recipient_string)

    response = deliver_mail

    expect(response.to).to eq [recipient_string]
    expect(response.cc).to eq []
    expect(response.bcc).to eq []
  end
end
```

Type `<Leader>s`:

```
rspec spec/recipient_interceptor_spec.rb:4
.

Finished in 0.03059 seconds
1 example, 0 failures
```

The screen is overtaken by a shell that runs only the focused spec.

Feeling good that this new spec passes,
run the whole file's specs with `<Leader>t`
to make sure the class's entire functionality is still intact:

```
rspec spec/recipient_interceptor_spec.rb
......

Finished in 0.17752 seconds
6 examples, 0 failures
```

Red, green, refactor.
From the program:

```ruby
def delivering_email(message)
  add_custom_headers message
  add_subject_prefix message
  message.to = @recipients
  message.cc = []
  message.bcc = []
end
```

Run `<Leader>l` without having to switch back to the spec:

```
rspec spec/recipient_interceptor_spec.rb
......

Finished in 0.17752 seconds
6 examples, 0 failures
```

Running specs in tight feedback loops
reduces switching cost between editor and shell,
making test-driven development easier.
