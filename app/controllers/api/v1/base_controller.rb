module Api
  module V1
    class BaseController < ApplicationController
      skip_before_action :require_authentication
      before_action :authenticate_api_token

      protect_from_forgery with: :null_session

      private
        def authenticate_api_token
          token = request.headers["Authorization"]&.delete_prefix("Bearer ")

          if token.present? && (user = User.active.find_by(api_token: token))
            Current.user = user
          else
            render json: { error: "Unauthorized" }, status: :unauthorized
          end
        end

        def ensure_editable(book)
          unless book.editable?
            render json: { error: "Forbidden" }, status: :forbidden
          end
        end

        def ensure_administrator
          unless Current.user.can_administer?
            render json: { error: "Forbidden — admin required" }, status: :forbidden
          end
        end
    end
  end
end
