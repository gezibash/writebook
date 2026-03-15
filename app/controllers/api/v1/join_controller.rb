module Api
  module V1
    class JoinController < BaseController
      skip_before_action :authenticate_api_token, only: :create

      def create
        unless Current.account.join_code == params[:join_code]
          render json: { error: "Invalid join code" }, status: :unauthorized
          return
        end

        user = User.create!(
          name: params[:name],
          email_address: params[:email],
          password: params[:password]
        )

        user.update!(api_token: SecureRandom.hex(32))
        render json: { token: user.api_token, user: { id: user.id, name: user.name, role: user.role } }, status: :created
      rescue ActiveRecord::RecordNotUnique
        render json: { error: "A user with that email already exists" }, status: :unprocessable_entity
      rescue ActiveRecord::RecordInvalid => e
        render json: { error: e.message }, status: :unprocessable_entity
      end
    end
  end
end
