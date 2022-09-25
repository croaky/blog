# begindoc: init
# config/initializers/workos.rb
WorkOS.configure do |config|
  config.key = ENV.fetch("WORKOS_API_KEY")
  config.timeout = 5
end
# enddoc: init

# begindoc: routes
# config/routes.rb
Rails.application.routes.draw do
  get "/login", to: "sso#login"
  get "/sso", to: "sso#callback"
  get "/logout", to: "sso#logout"
end
# enddoc: routes

# begindoc: controller
# app/controllers/sso_controller.rb
class SsoController < ApplicationController
  skip_before_action :require_login

  def login
    # Builds a string, does not make an HTTP request
    workos_auth_url = WorkOS::SSO.authorization_url(
      organization: ENV.fetch("WORKOS_ORGANIZATION"),
      client_id: ENV.fetch("WORKOS_CLIENT_ID"),
      redirect_uri: "#{base_url}/sso"
    )
    render locals: {
      workos_auth_url: workos_auth_url
    }
  end

  def callback
    cookies.clear

    if params["code"].nil?
      flash[:notice] = "Forbidden"
      redirect_to "/login"
    else
      begin
        # Makes an HTTP request, required for login to function
        profile = WorkOS::SSO.profile_and_token(
          code: params["code"],
          client_id: ENV.fetch("WORKOS_CLIENT_ID")
        ).profile
      rescue WorkOS::TimeoutError => err
        flash_msg = "Our SSO provider timed out. Check https://status.workos.com or wait a moment and try again."
        Sentry.capture_message("WorkOS::SSO.profile_and_token #{err}")
      end

      if profile
        user = User.active.by_normalized_email(profile.email)
      end

      if user
        remember(user)

        begin
          # Makes an HTTP request, but failure should not block login
          WorkOS::AuditTrail.create_event(
            event: {
              action_type: "read",
              action: "user.login_succeeded",
              actor_id: user.id.to_s,
              actor_name: user.full_name,
              group: "ivp.com",
              location: request.remote_ip,
              occurred_at: Time.current.iso8601,
              target_id: user.id.to_s,
              target_name: user.full_name
            }
          )
        rescue WorkOS::TimeoutError => err
          Sentry.capture_message(
            "WorkOS::AuditTrail.create_event #{err}",
            extra: {
              user_id: user.id.to_s,
              user_name: user.full_name
            }
          )
        end

        redirect_to session[:return_to] || "/"
      else
        flash[:error] = flash_msg || "Forbidden"
        redirect_to "/login"
      end
    end
  rescue => err
    Sentry.capture_exception(err)
    flash[:error] = "There was a problem. Please wait a moment and then try again."
    redirect_to "/login"
  end

  def logout
    forget_user
    flash[:notice] = "Logged out"
    redirect_to "/login"
  end
end
# enddoc: controller
