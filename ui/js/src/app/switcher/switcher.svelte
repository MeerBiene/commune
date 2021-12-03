<script>
import SwitcherHome from './home/home.svelte'
import SwitcherItems from './items/items.svelte'
import Footer from './footer/footer.svelte'
import Separator from './separator/separator.svelte'
import Tools from './tools/tools.svelte'
import { Router, Route } from 'svelte-navigator'

import { store } from '../store/store.js'

import Settings from './tools/settings/settings.svelte'

let mode = $store.settings.switcher.mode;

$: expanded = $store?.settings?.switcher?.mode == 'expanded'
$: collapsed = $store?.settings?.switcher?.mode == 'collapsed'

$: isMobile = $store.isMobile
$: mobileViewToggled = $store.mobileViewToggled

$: hideSwitcher = isMobile && !mobileViewToggled


</script>

<div class="switcher no-select"
    class:no-dis={hideSwitcher}
    class:switcher-collapsed={collapsed}
    class:switcher-expanded={expanded}>

    <SwitcherHome />



    <Separator />


    <Router primary={false}>

        <Route path="/:server/:room/*" let:params>
            <SwitcherItems 
            server={params.server}
            room={params.room}/>
        </Route>

        <Route let:params>
            <SwitcherItems />
        </Route>

    </Router>

    <Separator />

    <Tools />

    <Separator />

    <Settings />

    <Footer/>

</div>

<style>
.no-dis {
    display: none;
}
</style>
