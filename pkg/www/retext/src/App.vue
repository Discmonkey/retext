<template>
  <div id="app" class="container-fluid">
    <nav class="navbar navbar-expand-lg navbar-light bg-white border-bottom">
      <a class="navbar-brand" href="#">
        <img id="logo" src="./assets/qode.svg" alt="qode">
      </a>
      <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>

      <div class="collapse navbar-collapse" id="navbarSupportedContent">
        <ul class="navbar-nav mr-auto">
          <li class="nav-item">
            <router-link to="/upload" class="nav-link">Upload</router-link>
          </li>

          <li class="nav-item">
            <router-link to="/about" class="nav-link">About</router-link>
          </li>
          <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
              Dropdown
            </a>
            <div class="dropdown-menu" aria-labelledby="navbarDropdown">
              <a class="dropdown-item" href="#">Action</a>
              <a class="dropdown-item" href="#">Another action</a>
              <div class="dropdown-divider"></div>
              <a class="dropdown-item" href="#">Something else here</a>
            </div>
          </li>
          <li class="nav-item">
            <a class="nav-link disabled" href="#">Disabled</a>
          </li>
        </ul>

      </div>
    </nav>
      <nav>
          <breadcrumbs></breadcrumbs>
          <notifications id="notificationBox"/>
      </nav>
    <router-view/>
  </div>
</template>

<script>
import Breadcrumbs from "@/components/nav/Breadcrumbs";
import Notifications from "@/components/Notifications";
export default {
    components: {Breadcrumbs, Notifications},
}
</script>
<style>
:root {
    --bg-prim: white;
    --fg-prim: black;
    --button-prim: #0069d9;
}

@font-face {
    font-family: Roboto;
    src: url("assets/fonts/Quattrocento-Bold.ttf")
}
#app {
  font-family: Roboto, serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

#logo {
  height: 3em;
}

#notificationBox {
    position: absolute;
    right: 5px;
    bottom: 5px;
    z-index: 20;
}

/*#nav {*/
/*  padding: 30px;*/
/*}*/

/*#nav a {*/
/*  font-weight: bold;*/
/*  color: #2c3e50;*/
/*}*/

/*#nav a.router-link-exact-active {*/
/*  color: #42b983;*/
/*}*/

/*.container-fluid {*/
/*  padding: 5%;*/
/*}*/
</style>
