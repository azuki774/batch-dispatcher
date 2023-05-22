<!-- Todo.vue -->
<template>
    <div>
        <h3>{{ todo_title }}</h3>
        <table>
            <tr>
                <th>Index</th>
                <th>Title</th>
                <th>Important</th>
                <th>Completed</th>
            </tr>
            <tr v-for="(todo, index) in todos" v-bind:key="index">
                <td>{{ index }}</td>
                <td>{{ todo.title }}</td>
                <td>{{ todo.important ? 'Yes' : 'No' }}</td>
                <td><input type="checkbox" v-model="todo.completed"></td>
            </tr>
        </table>
        <hr />
        TODO:<input type="text" v-model="title">{{ title }}<br>
        IMPORTANT: <input type="checkbox" v-model="important"><br>
        <button @click="addTodo">Add New Todo</button>
    </div>
</template>
 
<script>
export default {
    props: { todo_title: String },
    data: function () {
        return {
            todos: [],
            title: '',
            important: false
        }
    },
    created: function () {
        console.log(this.todos);
    },
    methods: {
        addTodo() {
            if (!this.title) {
                alert("TODOに何も書いてません")
            }
            else {
                var new_todo = { title: this.title, important: this.important, completed: false }
                this.todos.push(new_todo);
                this.todos = this.todos.filter(todo => !todo.completed)
            }
        }
    }
}
</script>
 
<style scoped>
</style>
