/* eslint-disable */
<template>
  <div id="app">
    <router-view></router-view>
    <footer>
      <nav>
        <ul>
          <li>
            <router-link to="/" active-class="active">
              <Icon icon="tabler:home" width="24" />
              <span>首页</span>
            </router-link>
          </li>
          <li>
            <router-link to="/finance" active-class="active">
              <Icon icon="tabler:coin" width="24" />
              <span>理财</span>
            </router-link>
          </li>
          <li>
            <router-link to="/profile" active-class="active">
              <Icon icon="tabler:user" width="24" />
              <span>我的</span>
            </router-link>
          </li>
        </ul>
      </nav>
    </footer>
  </div>
</template>

<script>
import { Icon } from '@iconify/vue'

export default {
  name: 'App',
  components: {
    Icon
  }
}
</script>

<style>
#app {
  font-family: Arial, sans-serif;
  margin: 0;
  padding: 0;
  padding-bottom: 60px;
  /* text-align: center; */
}

footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #ffffff;
  padding: 8px 0;
  box-shadow: 0 -2px 10px rgba(0,0,0,0.05);
  z-index: 100;
  font-size: 0.9rem;
  color: #666;
}

footer nav {
  margin: 0;
}

footer nav ul {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  justify-content: space-around;
}

footer nav a {
  text-decoration: none;
  color: #666;
  display: flex;
  flex-direction: column;
  align-items: center;
  font-size: 12px;
  padding: 4px 0;
}

footer nav a.active {
  color: #2C5530;
}

footer nav a span {
  margin-top: 4px;
}
</style>