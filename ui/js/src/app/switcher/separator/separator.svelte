<script>
import { store } from '../../store/store.js'

$: expanded = $store.settings.switcher.mode == 'expanded'
$: collapsed = $store.settings.switcher.mode == 'collapsed'

</script>

<div class="separator-item">
    <div class="switcher-separator gr-center"
        class:collapsed={collapsed}
        class:expanded={expanded}>
    </div>
</div>


<style>
.separator-item {
    display: grid;
    border-radius: 50%;
    cursor: pointer;
}
.switcher-separator {
    width: 30px;
    height: 1px;
    border-radius: 1px;
    background-color: var(--separator);
}

.collapsed {
    width: 14px;
}
.expanded {
    width: 216px;
}
</style>
