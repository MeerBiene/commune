<script>
import { Router, Route } from 'svelte-navigator'
import Home from './home/home.svelte'
import Explore from './explore/explore.svelte'
import RoomContainer from './room/room-container.svelte'
import UserProfile from './user-profile/user-profile.svelte'
import UserDirectory from './user-directory/user-directory.svelte'
import { store } from '../store/store.js'



</script>


<Router primary={false}>
    <div class="view">

            <Home/>

            <Explore />

            <RoomContainer />

            <UserProfile />

            <UserDirectory />

    </div>
</Router>

<style>
.view {
}
</style>

