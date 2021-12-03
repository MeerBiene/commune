<script>
import { store } from '../../store/store.js'
import AddNew from './add-new/add-new.svelte'
import Explore from './explore/explore.svelte'
import UserDirectory from './user-directory/user-directory.svelte'

$: normal = $store.settings.switcher.mode == 'normal'
$: expanded = $store.settings.switcher.mode == 'expanded'
$: collapsed = $store.settings.switcher.mode == 'collapsed'

</script>

<div class="switcher-tools pv1" class:expanded={expanded}>

    <AddNew expanded={expanded} collapsed={collapsed}/>

    <Explore expanded={expanded} collapsed={collapsed}/>

    <UserDirectory expanded={expanded} collapsed={collapsed}/>

</div>

<style>
.switcher-tools {
    display: grid;
    grid-template-columns: auto;
    grid-template-rows: repeat(3,auto);
}

.expanded {
}
</style>
