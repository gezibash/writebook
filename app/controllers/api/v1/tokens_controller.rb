module Api
  module V1
    class TokensController < BaseController
      skip_before_action :authenticate_api_token, only: :create

      def create
        user = User.active.find_by(email_address: params[:email])

        if user&.authenticate(params[:password])
          user.update!(api_token: SecureRandom.hex(32)) if user.api_token.blank?
          render json: { token: user.api_token, user: { id: user.id, name: user.name, role: user.role } }
        else
          render json: { error: "Invalid credentials" }, status: :unauthorized
        end
      end
    end
  end
end
