<script>
import { useLocation } from "svelte-navigator";

const location = useLocation();

</script>

<div class="slct">
{$location.pathname}
</div>
