module Api
  module V1
    class SearchesController < BaseController
      def index
        query = params[:q]
        return render json: { error: "Missing query parameter: q" }, status: :unprocessable_entity if query.blank?

        results = []

        Book.accessable_or_published.ordered.each do |book|
          leaves = book.leaves.active.search(query).favoring_title.limit(50)
          next if leaves.empty?

          results << {
            book_id: book.id,
            book_title: book.title,
            results: leaves.map { |leaf| result_json(leaf, book) }
          }
        end

        render json: results
      end

      def show
        @book = Book.accessable_or_published.find(params[:book_id])
        query = params[:q]
        return render json: { error: "Missing query parameter: q" }, status: :unprocessable_entity if query.blank?

        leaves = @book.leaves.active.search(query).favoring_title.limit(50)
        render json: leaves.map { |leaf| result_json(leaf, @book) }
      end

      private
        def result_json(leaf, book)
          {
            id: leaf.id,
            title: leaf.title,
            title_snippet: strip_marks(leaf.title_match),
            content_snippet: strip_marks(leaf.content_match),
            type: leaf.leafable_type,
            book_id: book.id,
            book_title: book.title
          }
        end

        def strip_marks(text)
          text&.gsub(/<\/?mark>/, "")
        end
    end
  end
end
