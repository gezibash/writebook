require "test_helper"

class Api::V1::SearchesControllerTest < ActionDispatch::IntegrationTest
  setup do
    @user = users(:david)
    @user.update!(api_token: "test-token-david")
    @headers = { "Authorization" => "Bearer test-token-david" }
    @book = books(:handbook)
  end

  # Global search

  test "global search returns results grouped by book" do
    post api_v1_book_pages_url(@book), headers: @headers, as: :json,
      params: { leaf: { title: "Deployment Workflow" }, page: { body: "The deployment process starts with staging" } }
    assert_response :created

    get api_v1_search_url(q: "deployment"), headers: @headers, as: :json
    assert_response :success

    json = JSON.parse(response.body)
    assert json.is_a?(Array)
    group = json.find { |g| g["book_id"] == @book.id }
    assert group, "Expected results grouped under book #{@book.id}"
    assert group["book_title"].present?
    assert group["results"].any? { |r| r["title"] == "Deployment Workflow" }
  end

  test "global search returns empty array for no matches" do
    get api_v1_search_url(q: "xyznonexistent"), headers: @headers, as: :json
    assert_response :success

    json = JSON.parse(response.body)
    assert_equal [], json
  end

  test "global search returns 422 for missing q param" do
    get api_v1_search_url, headers: @headers, as: :json
    assert_response :unprocessable_entity
  end

  test "global search rejects unauthenticated requests" do
    get api_v1_search_url(q: "test"), as: :json
    assert_response :unauthorized
  end

  # Book-scoped search

  test "book-scoped search returns results for the book" do
    post api_v1_book_pages_url(@book), headers: @headers, as: :json,
      params: { leaf: { title: "Monitoring Setup" }, page: { body: "Configure monitoring alerts and dashboards" } }
    assert_response :created

    get api_v1_book_search_url(@book, q: "monitoring"), headers: @headers, as: :json
    assert_response :success

    json = JSON.parse(response.body)
    assert json.is_a?(Array)
    assert json.any? { |r| r["title"] == "Monitoring Setup" }
  end

  test "book-scoped search returns empty array for no matches" do
    get api_v1_book_search_url(@book, q: "xyznonexistent"), headers: @headers, as: :json
    assert_response :success

    json = JSON.parse(response.body)
    assert_equal [], json
  end

  test "book-scoped search returns 422 for missing q param" do
    get api_v1_book_search_url(@book), headers: @headers, as: :json
    assert_response :unprocessable_entity
  end

  test "book-scoped search rejects unauthenticated requests" do
    get api_v1_book_search_url(@book, q: "test"), as: :json
    assert_response :unauthorized
  end

  # Result format

  test "search results contain no mark HTML tags" do
    post api_v1_book_pages_url(@book), headers: @headers, as: :json,
      params: { leaf: { title: "Tagged Content" }, page: { body: "This has searchable tagged content inside" } }
    assert_response :created

    get api_v1_book_search_url(@book, q: "tagged"), headers: @headers, as: :json
    assert_response :success

    json = JSON.parse(response.body)
    assert json.any?, "Expected at least one result"
    json.each do |result|
      assert_no_match(/<\/?mark>/, result["title_snippet"].to_s)
      assert_no_match(/<\/?mark>/, result["content_snippet"].to_s)
    end
  end
end
