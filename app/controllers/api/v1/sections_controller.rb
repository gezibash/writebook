module Api
  module V1
    class SectionsController < BaseController
      before_action :set_book
      before_action -> { ensure_editable(@book) }, except: %i[ index ]
      before_action :set_leaf, only: %i[ show update destroy ]

      def index
        sections = @book.leaves.active.with_leafables.positioned.where(leafable_type: "Section")
        render json: sections.map { |l| section_json(l) }
      end

      def show
        render json: section_json(@leaf)
      end

      def create
        section = Section.new(section_params)
        leaf = @book.press(section, leaf_params)

        if position = params[:position]&.to_i
          leaf.move_to_position(position)
        end

        render json: section_json(leaf), status: :created
      end

      def update
        @leaf.edit(leafable_params: section_params, leaf_params: leaf_params)
        render json: section_json(@leaf.reload)
      end

      def destroy
        @leaf.trashed!
        head :no_content
      end

      private
        def set_book
          @book = Book.accessable_or_published.find(params[:book_id])
        end

        def set_leaf
          @leaf = @book.leaves.active.find(params[:id])
        end

        def leaf_params
          title = params.dig(:leaf, :title) || params.dig(:section, :body).to_s.truncate(50) || "Section"
          { title: title }
        end

        def section_params
          params.fetch(:section, {}).permit(:body, :theme)
        end

        def section_json(leaf)
          {
            id: leaf.id,
            title: leaf.title,
            type: "Section",
            body: leaf.leafable.body,
            theme: leaf.leafable.theme,
            status: leaf.status,
            position: leaf.position_score,
            slug: leaf.slug,
            created_at: leaf.created_at,
            updated_at: leaf.updated_at
          }
        end
    end
  end
end
