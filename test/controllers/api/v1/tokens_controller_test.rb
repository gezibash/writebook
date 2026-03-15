require "test_helper"

class Api::V1::TokensControllerTest < ActionDispatch::IntegrationTest
  test "create returns api token with valid credentials" do
    post api_v1_tokens_url, params: { email: "david@37signals.com", password: "secret123456" }, as: :json

    assert_response :success
    json = JSON.parse(response.body)
    assert json["token"].present?
    assert_equal "David", json["user"]["name"]
    assert_equal "administrator", json["user"]["role"]
  end

  test "create rejects invalid credentials" do
    post api_v1_tokens_url, params: { email: "david@37signals.com", password: "wrong" }, as: :json

    assert_response :unauthorized
  end

  test "create returns same token on repeated calls" do
    post api_v1_tokens_url, params: { email: "david@37signals.com", password: "secret123456" }, as: :json
    token1 = JSON.parse(response.body)["token"]

    post api_v1_tokens_url, params: { email: "david@37signals.com", password: "secret123456" }, as: :json
    token2 = JSON.parse(response.body)["token"]

    assert_equal token1, token2
  end
end
