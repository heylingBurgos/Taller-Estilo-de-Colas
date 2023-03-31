<template>
  <tr>
    <th scope="row">{{item.id}}</th>
    <td>{{item.nombre}}</td>
    <td>
        <button class="btn btn-info btn-sm">Consumir</button>
    </td>
  </tr>
</template>

<script>
export default {
    props: ['item']
}
</script>

<style>

</style>