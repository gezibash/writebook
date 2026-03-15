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
    cover_seed = "feedfacecafebeef"

    assert_difference -> { Book.count }, +1 do
      post api_v1_books_url, headers: @headers, as: :json,
        params: { book: { title: "API Book", subtitle: "Built via API", author: "Agent", cover_style: "glass", cover_seed: cover_seed } }
    end

    assert_response :created
    json = JSON.parse(response.body)
    assert_equal "API Book", json["title"]
    assert_equal "Agent", json["author"]
    assert_equal "glass", json["cover_style"]
    assert_equal cover_seed, json["cover_seed"]
  end

  test "update modifies a book" do
    book = books(:handbook)
    patch api_v1_book_url(book), headers: @headers, as: :json,
      params: { book: { title: "Updated Handbook", cover_style: "rings", cover_seed: "0123456789abcdef" } }

    assert_response :success
    assert_equal "Updated Handbook", book.reload.title
    assert_equal "rings", book.cover_style
    assert_equal "0123456789abcdef", book.cover_seed
  end

  test "update toggles published status" do
    book = books(:handbook)
    patch api_v1_book_url(book), headers: @headers, as: :json,
      params: { book: { published: true } }

    assert_response :success
    assert book.reload.published?
  end

  test "destroy deletes a book" do
    book = books(:handbook)
    assert_difference -> { Book.count }, -1 do
      delete api_v1_book_url(book), headers: @headers, as: :json
    end

    assert_response :no_content
  end
end
