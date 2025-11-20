<template>
  <div class="articles">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">文章列表</h1>
        <p class="page-subtitle">探索知识的海洋</p>
      </div>

      <!-- 筛选和排序 -->
      <div class="filters">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索文章..."
          clearable
          style="width: 300px"
          @change="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>

        <div class="filter-actions">
          <el-select
            v-model="selectedCategory"
            placeholder="选择分类"
            clearable
            style="width: 150px"
            @change="handleFilter"
          >
            <el-option
              v-for="category in categories"
              :key="category"
              :label="category"
              :value="category"
            />
          </el-select>

          <el-select
            v-model="sortBy"
            placeholder="排序方式"
            style="width: 150px"
            @change="handleFilter"
          >
            <el-option label="最新发布" value="created_at" />
            <el-option label="最多浏览" value="views" />
            <el-option label="最多点赞" value="likes" />
          </el-select>
        </div>
      </div>

      <!-- 文章列表 -->
      <div v-loading="loading" class="articles-list">
        <ArticleCard
          v-for="article in articles"
          :key="article.id"
          :article="article"
        />

        <el-empty v-if="!loading && !articles.length" description="暂无文章" />
      </div>

      <!-- 分页 -->
      <div v-if="total > pageSize" class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next, jumper"
          @current-change="handlePageChange"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getArticles, searchArticles } from '@/api/article'
import ArticleCard from '@/components/ArticleCard.vue'

const route = useRoute()

const articles = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const searchKeyword = ref('')
const selectedCategory = ref('')
const sortBy = ref('created_at')
const categories = ref(['前端开发', '后端开发', '数据库', '运维部署', '算法', '其他'])

onMounted(() => {
  // 从URL获取查询参数
  searchKeyword.value = route.query.keyword || ''
  selectedCategory.value = route.query.category || ''

  fetchArticles()
})

watch(() => route.query, () => {
  searchKeyword.value = route.query.keyword || ''
  selectedCategory.value = route.query.category || ''
  currentPage.value = 1
  fetchArticles()
})

const fetchArticles = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      page_size: pageSize.value,
      sort: sortBy.value
    }

    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }

    if (selectedCategory.value) {
      params.category = selectedCategory.value
    }

    const apiCall = searchKeyword.value ? searchArticles : getArticles
    const { data } = await apiCall(params)

    articles.value = data.list || []
    total.value = data.total || 0
  } catch (error) {
    console.error('Failed to fetch articles:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  currentPage.value = 1
  fetchArticles()
}

const handleFilter = () => {
  currentPage.value = 1
  fetchArticles()
}

const handlePageChange = () => {
  fetchArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<style scoped>
.articles {
  padding: 20px 0;
}

.filters {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.filter-actions {
  display: flex;
  gap: 12px;
}

.articles-list {
  min-height: 500px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}

@media (max-width: 768px) {
  .filters {
    flex-direction: column;
    gap: 16px;
  }

  .filter-actions {
    width: 100%;
    flex-direction: column;
  }

  .filter-actions .el-select {
    width: 100% !important;
  }
}
</style>
