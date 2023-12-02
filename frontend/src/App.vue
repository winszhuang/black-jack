<script setup lang="ts">
</script>

<template>
  <h1 class="title">BLACK JACK</h1>
    <div class="main">
        <h2><span id="command">Gambling Time</span></h2>
        <div class="row1">
            <div id="dealer-box">
                <h2>Dealer: <span id="dealerscore">0</span></h2>
            </div>
            <div id="ihjiljiljil">
                <h2>You: 
                    <span id="yourscore">0</span>
                    <div id="your-box" class="flex flex-wrap">

                    </div>
                </h2>
            </div>
            <div id="test">
                <h2>You: <span id="testscore">0</span></h2>
            </div>
        </div>
        <div class="row2">
            <div class="buttons">
                <button class="btn-lg btn-danger" id="deal" onmouseover="tink.play()">Ready</button>
                <button class="btn-lg btn-success" id="hit" onmouseover="tink.play()" > Hit </button> 
                <button class="btn-lg btn-warning" id="stand" onmouseover="tink.play()">Stand</button>
            </div>
        </div>
        <div class="row3">
            <h3>Scoreboard: </h3>
            <table>
                <tr>
                    <th>Wins</th>
                    <th>Losses</th>
                    <th>Draws</th>
                </tr>
                <tr>
                    <td><span id="wins">0</span></td>
                    <td><span id="losses">0</span></td>
                    <td><span id="draws">0</span></td>
                </tr>
            </table>
        </div>
    </div>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}
.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}
.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
