require "test_helper"

class Api::V1::JoinControllerTest < ActionDispatch::IntegrationTest
  test "join creates a user and returns a token" do
    join_code = accounts(:signal).join_code

    assert_difference -> { User.count }, +1 do
      post api_v1_join_url, as: :json,
        params: { join_code: join_code, name: "New Agent", email: "agent@example.com", password: "secret123" }
    end

    assert_response :created
    json = JSON.parse(response.body)
    assert json["token"].present?
    assert_equal "New Agent", json["user"]["name"]
  end

  test "join rejects invalid join code" do
    post api_v1_join_url, as: :json,
      params: { join_code: "bad-code", name: "Nope", email: "nope@example.com", password: "secret123" }

    assert_response :unauthorized
    json = JSON.parse(response.body)
    assert_equal "Invalid join code", json["error"]
  end

  test "join rejects duplicate email" do
    join_code = accounts(:signal).join_code

    post api_v1_join_url, as: :json,
      params: { join_code: join_code, name: "Dup", email: users(:david).email_address, password: "secret123" }

    assert_response :unprocessable_entity
  end
end
