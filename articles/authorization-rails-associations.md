# Authorization with Rails Associations

Authorize requests through the domain model's associations:

```ruby
class User < ApplicationRecord
  has_many :brands
end

class Brand < ApplicationRecord
  has_many :offers
end

Rails.application.routes.draw do
  resources :brands, only: [] do
    resources :offers, only: [:new]
  end
end
```

An example controller test using [Clearance]'s `sign_in_as`:

[Clearance]: https://github.com/thoughtbot/clearance

```ruby
it "does not find brands unassociated with user" do
  sign_in_as create(:user)

  expect { get :new, brand_id: 1 }.to raise_error(ActiveRecord::RecordNotFound }
end
```

`ActiveRecord::RecordNotFound` is raised because
there is no record of the user with this brand,
returning a 404.

Make the tests pass by restricting users to their brands
using [Clearance]'s `:require_login` and `current_user`:

```ruby
class OffersController < ApplicationController
  before_filter :require_login

  def new
    @brand = current_user.brands.find(params[:brand_id])
    @offer = @brand.offers.build
  end
end
```

This authorization approach requires few lines of code
and no extra dependencies beyond Rails and Clearance.
It removes duplication,
is testable,
and uses conventions.
