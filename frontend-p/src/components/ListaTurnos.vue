<template>
  <div class="my-5">
    <h4>Lista de turnos</h4>
    <table class="table">
        <thead>
            <tr>
                <th scope="col">#</th>
                <th scope="col">Nombre</th>
                <th scope="col">Accion</th>
            </tr>
        </thead>
        <tbody id="items">
            <Item v-for="item in items" v-bind:key="item.cedula" :item="item">
            </Item>
            />
        </tbody>
        <tfoot>
            <tr id="footer-lista">
                <th scope="row" colspan="5">Lista vacia - esperando turnos</th>
            </tr>
        </tfoot>
    </table>
  </div>
</template>

<script>
import { useStore } from 'vuex'
import { computed } from 'vue'
import Item from './Item.vue'

export default {
    components: {
        Item
    },
    setup() {
        const store = useStore()
        const items = computed(() => store.state.turnos)

        return {
            items
        }
    }

}
</script>

<style>

</style>