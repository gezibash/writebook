require "test_helper"

class Api::V1::BooksControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = users(:david)
    @user.update!(api_token: "test-token-david")
    @headers = { "Authorization" => "Bearer test-token-david" }
  end

  test "index returns books the user can access" do
    get api_v1_books_url, headers: @headers, as: :json

    assert_response :success
    json = JSON.parse(response.body)
    assert json.is_a?(Array)
    assert json.any? { |b| b["title"] == "Handbook" }
  end

  test "index rejects unauthenticated requests" do
    get api_v1_books_url, as: :json

    assert_response :unauthorized
  end

  test "show returns book with leaves" do
    book = books(:handbook)
    get api_v1_book_url(book), headers: @headers, as: :json

    assert_response :success
    json = JSON.parse(response.body)
    assert_equal "Handbook", json["title"]
    assert json["leaves"].is_a?(Array)
  end

  test "create makes a new book" do
    assert_difference -> { Book.count }, +1 do
      post api_v1_books_url, headers: @headers, as: :json,
        params: { book: { title: "API Book", subtitle: "Built via API", author: "Agent" } }
    end

    assert_response :created
    json = JSON.parse(response.body)
    assert_equal "API Book", json["title"]
    assert_equal "Agent", json["author"]
  end

  test "update modifies a book" do
    book = books(:handbook)
    patch api_v1_book_url(book), headers: @headers, as: :json,
      params: { book: { title: "Updated Handbook" } }

    assert_response :success
    assert_equal "Updated Handbook", book.reload.title
  end

  test "destroy deletes a book" do
    book = books(:handbook)
    assert_difference -> { Book.count }, -1 do
      delete api_v1_book_url(book), headers: @headers, as: :json
    end

    assert_response :no_content
  end
end
