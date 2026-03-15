require "test_helper"

class Api::V1::PagesControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = users(:david)
    @user.update!(api_token: "test-token-david")
    @headers = { "Authorization" => "Bearer test-token-david" }
    @book = books(:handbook)
  end

  test "index lists pages in a book" do
    post api_v1_book_pages_url(@book), headers: @headers, as: :json,
      params: { leaf: { title: "Test Page" }, page: { body: "content" } }

    get api_v1_book_pages_url(@book), headers: @headers, as: :json

    assert_response :success
    json = JSON.parse(response.body)
    assert json.is_a?(Array)
    assert json.any? { |p| p["title"] == "Test Page" }
  end

  test "create adds a page to a book" do
    assert_difference -> { @book.leaves.count }, +1 do
      post api_v1_book_pages_url(@book), headers: @headers, as: :json,
        params: { leaf: { title: "Chapter 1" }, page: { body: "# Hello\n\nThis is markdown." } }
    end

    assert_response :created
    json = JSON.parse(response.body)
    assert_equal "Chapter 1", json["title"]
    assert_equal "Page", json["type"]
    assert_includes json["body"], "Hello"
  end

  test "update modifies page content" do
    post api_v1_book_pages_url(@book), headers: @headers, as: :json,
      params: { leaf: { title: "Draft" }, page: { body: "initial" } }
    leaf_id = JSON.parse(response.body)["id"]

    patch api_v1_book_page_url(@book, leaf_id), headers: @headers, as: :json,
      params: { leaf: { title: "Final Chapter" }, page: { body: "updated content" } }

    assert_response :success
    json = JSON.parse(response.body)
    assert_equal "Final Chapter", json["title"]
  end

  test "destroy trashes a page" do
    post api_v1_book_pages_url(@book), headers: @headers, as: :json,
      params: { leaf: { title: "To Delete" }, page: { body: "bye" } }
    leaf_id = JSON.parse(response.body)["id"]

    delete api_v1_book_page_url(@book, leaf_id), headers: @headers, as: :json

    assert_response :no_content
    assert_equal "trashed", Leaf.find(leaf_id).status
  end

  test "rejects unauthenticated requests" do
    post api_v1_book_pages_url(@book), as: :json,
      params: { leaf: { title: "Nope" }, page: { body: "nope" } }

    assert_response :unauthorized
  end
end
