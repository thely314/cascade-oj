<script setup lang="ts">import { ref, computed, onMounted } from 'vue';
import { getContests, getSingleContest, postContest, putContest, deleteContest, getProblems, getContestCompetitors, putContestCompetitors, getUsers } from '../../api/admin';
import type { ContestMetadata, GetSingleContestReply, PostContestRequest, PutContestRequest, ProblemMetadata, UserInfo, GetContestCompetitorsReply, PutContestCompetitorsRequest } from '../../api/types';
const contests = ref<ContestMetadata[]>([]);
const loading = ref(false);
const error = ref('');
const showCreateModal = ref(false);
const showEditModal = ref(false);
const showProblemModal = ref(false);
const showCompetitorModal = ref(false);
const currentContest = ref<ContestMetadata | null>(null);
const selectedProblems = ref<number[]>([]);
const allProblems = ref<ProblemMetadata[]>([]);
const problemSearchQuery = ref('');
const selectedUsers = ref<number[]>([]);
const allUsers = ref<UserInfo[]>([]);
const userSearchQuery = ref('');
const form = ref({
 title: '',
 description: '',
 startTime: '',
 endTime: ''
});
const statusTone: Record<string, string> = {
 upcoming: 'badge-muted',
 ongoing: 'badge-live',
 ended: 'badge-dim'
};
const filteredProblems = computed(() => {
	if (!problemSearchQuery.value)
		return allProblems.value;
	const query = problemSearchQuery.value.toLowerCase();
	return allProblems.value.filter(p => 
		p.title.toLowerCase().includes(query) || 
		(p.description && p.description.toLowerCase().includes(query))
	);
});
const filteredUsers = computed(() => {
	if (!userSearchQuery.value)
		return allUsers.value;
	const query = userSearchQuery.value.toLowerCase();
	return allUsers.value.filter(u => u.username.toLowerCase().includes(query) || u.email.toLowerCase().includes(query));
});
const formatDate = (dateStr: string) => {
  if (!dateStr)
    return '';
  return new Date(dateStr).toLocaleString('zh-CN');
};
const formatDateTimeLocal = (dateStr: string) => {
  if (!dateStr)
    return '';
  const date = new Date(dateStr);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  return `${year}-${month}-${day}T${hours}:${minutes}`;
};
const fetchContests = async () => {
 loading.value = true;
 error.value = '';
 try {
 const response = await getContests();
 contests.value = response.contests || [];
 }
 catch (err: any) {
 error.value = 'Failed to load contests';
 console.error(err);
 }
 finally {
 loading.value = false;
 }
};
const fetchProblems = async () => {
	try {
		const response = await getProblems();
		allProblems.value = response.problems || [];
	}
	catch (err) {
		console.error('Failed to load problems:', err);
	}
};
const fetchUsers = async () => {
	try {
		const response = await getUsers();
		allUsers.value = response.users || [];
	}
	catch (err) {
		console.error('Failed to load users:', err);
	}
};
const openCreateModal = () => {
  const now = new Date();
  const oneHourLater = new Date(now.getTime() + 60 * 60 * 1000);
  const formatDateTimeLocal = (date: Date) => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  };
  form.value = {
    title: '',
    description: '',
    startTime: formatDateTimeLocal(now),
    endTime: formatDateTimeLocal(oneHourLater)
  };
  selectedProblems.value = [];
  showCreateModal.value = true;
};
const openEditModal = async (contestId: number) => {
  try {
    const response = await getSingleContest(contestId);
    console.log('Response from getSingleContest:', response);
    const problemIds = response.problemIds || [];
    console.log('Problem IDs:', problemIds);
    currentContest.value = response.metadata;
    form.value = {
      title: response.metadata.title,
      description: response.description,
      startTime: formatDateTimeLocal(response.metadata.startTime),
      endTime: formatDateTimeLocal(response.metadata.endTime)
    };
    selectedProblems.value = [...problemIds];
    console.log('Selected problems after assignment:', selectedProblems.value);
    showEditModal.value = true;
    console.log('showEditModal set to:', showEditModal.value);
  }
  catch (err) {
    console.error('Failed to load contest:', err);
  }
};
const openProblemModal = () => {
	fetchProblems();
	showProblemModal.value = true;
};
const openCompetitorModal = async (contestId: number) => {
	currentContest.value = { id: contestId, title: '', startTime: '', endTime: '', status: '' };
	try {
		const response: GetContestCompetitorsReply = await getContestCompetitors(contestId);
		selectedUsers.value = response.competitors.map(c => c.userId);
	}
	catch (err) {
		console.error('Failed to load competitors:', err);
		selectedUsers.value = [];
	}
	fetchUsers();
	showCompetitorModal.value = true;
};
const toggleProblemSelection = (problemId: number) => {
 const index = selectedProblems.value.indexOf(problemId);
 if (index === -1) {
 selectedProblems.value.push(problemId);
 }
 else {
 selectedProblems.value.splice(index, 1);
 }
};
const toggleUserSelection = (userId: number) => {
	const index = selectedUsers.value.indexOf(userId);
	if (index === -1) {
		selectedUsers.value.push(userId);
	}
	else {
		selectedUsers.value.splice(index, 1);
	}
};
const handleSaveCompetitors = async () => {
	if (!currentContest.value)
		return;
	try {
		const data: PutContestCompetitorsRequest = {
			userIds: selectedUsers.value
		};
		await putContestCompetitors(currentContest.value.id, data);
		alert('参赛人员更新成功！');
		showCompetitorModal.value = false;
	}
	catch (err: any) {
		let errorMessage = '操作失败，请稍后重试';
		if (err.response?.data?.message) {
			errorMessage = err.response.data.message;
		}
		else if (err.message) {
			errorMessage = err.message;
		}
		alert(`更新失败: ${errorMessage}`);
	}
};
const handleSubmit = async (isEdit: boolean) => {
  if (!form.value.title) {
    alert('请输入比赛标题');
    return;
  }
  if (!form.value.startTime || !form.value.endTime) {
    alert('请选择开始时间和结束时间');
    return;
  }
  if (selectedProblems.value.length === 0) {
    alert('请至少选择一个题目');
    return;
  }
  if (new Date(form.value.startTime) >= new Date(form.value.endTime)) {
    alert('结束时间必须晚于开始时间');
    return;
  }
  try {
    const data: PostContestRequest | PutContestRequest = {
      title: form.value.title,
      description: form.value.description,
      startTime: new Date(form.value.startTime).toISOString(),
      endTime: new Date(form.value.endTime).toISOString(),
      problems: {
        problemIds: selectedProblems.value
      }
    };
    if (isEdit && currentContest.value) {
      await putContest(currentContest.value.id, data as PutContestRequest);
      alert(`比赛 "${form.value.title}" 修改成功！`);
      showEditModal.value = false;
      fetchContests();
    }
    else {
      await postContest(data as PostContestRequest);
      alert(`比赛 "${form.value.title}" 创建成功！`);
      showCreateModal.value = false;
      fetchContests();
    }
  }
  catch (err: any) {
    let errorMessage = '操作失败，请稍后重试';
    if (err.response?.data?.message) {
      errorMessage = err.response.data.message;
    } else if (err.message) {
      errorMessage = err.message;
    } else if (err.response?.data?.error) {
      errorMessage = err.response.data.error;
    }
    alert(`创建/修改失败: ${errorMessage}\n\n请检查表单内容并重试。`);
  }
};
const handleDelete = async (contestId: number) => {
 if (!confirm('Are you sure you want to delete this contest?')) {
 return;
 }
 try {
 await deleteContest(contestId);
 alert('Contest deleted successfully');
 fetchContests();
 }
 catch (err: any) {
 const errorMessage = err.response?.data?.message || 'Delete failed';
 alert(`Failed: ${errorMessage}`);
 }
};
onMounted(() => {
 fetchContests();
});
</script>

<template>
  <div class="page">
    <header class="page-header">
      <div>
        <p class="eyebrow">Admin</p>
        <h1 class="title">Contests</h1>
      </div>
      <div class="actions">
        <button class="ghost" @click="fetchContests" :disabled="loading">
          {{ loading ? 'Loading...' : 'Refresh' }}
        </button>
        <button class="primary" @click="openCreateModal">New Contest</button>
      </div>
    </header>

    <section class="panel">
      <header class="panel-header">
        <h2>Contest List</h2>
        <a href="#">View all</a>
      </header>
      <div v-if="error" class="error-message">{{ error }}</div>
      <ul v-else class="list">
        <li v-for="contest in contests" :key="contest.id" class="list-item">
          <div class="list-main">
            <p class="list-title">{{ contest.title }}</p>
            <p class="list-meta">ID {{ contest.id }} · {{ formatDate(contest.startTime) }} → {{ formatDate(contest.endTime) }}</p>
          </div>
          <div class="list-right">
            <span class="badge" :class="statusTone[contest.status] || 'badge-muted'">{{ contest.status }}</span>
            <button class="ghost" @click="openCompetitorModal(contest.id)">Manage Competitors</button>
            <button class="ghost" @click="openEditModal(contest.id)">Edit</button>
            <button class="ghost" @click="handleDelete(contest.id)">Delete</button>
          </div>
        </li>
      </ul>
    </section>

    <div v-if="showCreateModal || showEditModal" class="modal-overlay" @click.self="showCreateModal = showEditModal = false">
      <div class="modal">
        <header class="modal-header">
          <h3>{{ showEditModal ? 'Edit Contest' : 'Create New Contest' }}</h3>
          <button class="close" @click="showCreateModal = showEditModal = false">×</button>
        </header>
        <div class="modal-body">
          <div class="form-group">
            <label>Title</label>
            <input v-model="form.title" type="text" placeholder="Enter contest title" />
          </div>
          <div class="form-group">
            <label>Description</label>
            <textarea v-model="form.description" placeholder="Enter contest description"></textarea>
          </div>
          <div class="form-group">
            <label>Start Time</label>
            <input v-model="form.startTime" type="datetime-local" />
          </div>
          <div class="form-group">
            <label>End Time</label>
            <input v-model="form.endTime" type="datetime-local" />
          </div>
          <div class="form-group">
            <label>Selected Problems ({{ selectedProblems.length }})</label>
            <button class="ghost" @click="openProblemModal">Select Problems</button>
            <div v-if="selectedProblems.length > 0" class="selected-problems">
              <span v-for="pid in selectedProblems" :key="pid" class="tag">{{ pid }}</span>
            </div>
          </div>
        </div>
        <footer class="modal-footer">
          <button class="ghost" @click="showCreateModal = showEditModal = false">Cancel</button>
          <button class="primary" @click="handleSubmit(showEditModal)">
            {{ showEditModal ? 'Update' : 'Create' }}
          </button>
        </footer>
      </div>
    </div>

    <div v-if="showProblemModal" class="modal-overlay" @click.self="showProblemModal = false">
      <div class="modal problem-modal">
        <header class="modal-header">
          <h3>Select Problems</h3>
          <button class="close" @click="showProblemModal = false">×</button>
        </header>
        <div class="modal-body">
          <div class="search-box">
            <input v-model="problemSearchQuery" type="text" placeholder="Search by title or description..." />
          </div>
          <table class="problem-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Title</th>
                <th>Description</th>
                <th>Select</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="problem in filteredProblems" :key="problem.id">
                <td>{{ problem.id }}</td>
                <td>{{ problem.title }}</td>
                <td>{{ problem.description }}</td>
                <td>
                  <input
                    type="checkbox"
                    :checked="selectedProblems.includes(problem.id)"
                    @change="toggleProblemSelection(problem.id)"
                  />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <footer class="modal-footer">
          <button class="ghost" @click="showProblemModal = false">Cancel</button>
          <button class="primary" @click="showProblemModal = false">Confirm Selection</button>
        </footer>
      </div>
    </div>

    <div v-if="showCompetitorModal" class="modal-overlay" @click.self="showCompetitorModal = false">
      <div class="modal problem-modal">
        <header class="modal-header">
          <h3>Manage Competitors</h3>
          <button class="close" @click="showCompetitorModal = false">×</button>
        </header>
        <div class="modal-body">
          <div class="form-group">
            <label>Selected Competitors ({{ selectedUsers.length }})</label>
            <div v-if="selectedUsers.length > 0" class="selected-problems">
              <span v-for="uid in selectedUsers" :key="uid" class="tag">{{ uid }}</span>
            </div>
          </div>
          <div class="search-box">
            <input v-model="userSearchQuery" type="text" placeholder="Search by username or email..." />
          </div>
          <table class="problem-table">
            <thead>
              <tr>
                <th>ID</th>
                <th>Username</th>
                <th>Email</th>
                <th>Select</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="user in filteredUsers" :key="user.userId">
                <td>{{ user.userId }}</td>
                <td>{{ user.username }}</td>
                <td>{{ user.email }}</td>
                <td>
                  <input
                    type="checkbox"
                    :checked="selectedUsers.includes(user.userId)"
                    @change="toggleUserSelection(user.userId)"
                  />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <footer class="modal-footer">
          <button class="ghost" @click="showCompetitorModal = false">Cancel</button>
          <button class="primary" @click="handleSaveCompetitors">Save</button>
        </footer>
      </div>
    </div>
  </div>
</template>

<style scoped src="./Contest.css"></style>