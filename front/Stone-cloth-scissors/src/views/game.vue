<template>
  <div class="game">
    <div class="head_title">
      <h1><b>{{ $t('roomId') }}:2131231{{ roomId }}</b></h1>
    </div>
    <div class="main">
      <div class="statistics">
        <div class="blico win">{{ $t('win') }}: <span>{{ win }}</span></div>
        <div class="blico draw">{{ $t('draw') }}:<span>{{ draw }}</span></div>
        <div class="blico lose">{{ $t('lose') }}:<span>{{ lose }}</span></div>
      </div>
      <div class="competition">
        <div class="block own_venue">
          <div class="frame">
            <div class="fist_area">
              <img v-if="ownImgIndex != null" class="img" :src="imgList[ownImgIndex]" alt="加载失败...">
              <div v-else class="text">?</div>
            </div>
            <div class="name">
              {{ $t('gamePlayer_you') }}
            </div>
          </div>
        </div>
        <div class="block VS">
          VS
        </div>
        <div class="block opponent_venue">
          <div class="frame">
            <div class="fist_area">
              <img v-if="opponentImgIndex != null" class="img" :src="imgList[opponentImgIndex]" alt="加载失败...">
              <div v-else class="text">?</div>
            </div>
            <div class="name">
              {{ $t('gamePlayer_opponent') }}
            </div>
          </div>
        </div>
      </div>
      <div class="fill_blcok"></div>
    </div>
    <div class="skill_venue">
      <div v-for="(item, i) in imgList" :key="i" :class="'skill' + i" @click="onClick(i + 1)">
        <img class="img" :src="item" alt="加载失败...">
      </div>
    </div>
    <div class="History" ref="overflow">
      <div v-for="(item,i) in HistoryList" :key="i" style="margin-bottom:1.25rem">
        <div class="title" v-if="language == 'zh'">
          第 {{ i+1 }} 局
        </div>
        <div class="title" v-else>
          {{ i+1 }}st round
        </div>
        <div style="padding-left:0.625rem">
          <div>{{ item.time }}</div>
          <div>{{ item.youHandNum }}</div> 
          <div>{{ item.opponentHandNum }}</div>
          <div>{{ item.isWin }}</div>
        </div>
      </div>
    </div>
    <Popup :show="startGameShow">
      <div class="Popup">
        {{ startGame }}
      </div>
    </Popup>
    <Popup :show="waitingShow">
      <div class="Popup">
        {{ waiting }}
      </div>
    </Popup>
    <Popup :show="resultShow">
      <div class="result_Popup">
        <div class="img_block">
          <img class="img" :src="resultImg" alt="加载失败...">
        </div>
        <div class="btnList">
          <el-button class="btn continue" @click="continueFun">{{ $t('continue') }}</el-button>
          <el-button class="btn exit" @click="exit">{{ $t('exit') }}</el-button>
        </div>
      </div>
    </Popup>
    <audio src="../assets/bj.mp3" :loop="true" ref="music" autoplay hidden :muted="muted"></audio>
  </div>
</template>

<script>
import store from "@/store";
import { mapGetters } from 'vuex'
import { sendSock } from "../utils/socket";
import Popup from '@/components/mask/mask.vue';
import { getDefaultLang,formatterDate } from "../utils/utils";
import Cookies from 'js-cookie';
import zh_win from '../assets/win.png';
import en_win from '../assets/Winning.png';
import zh_lose from '../assets/lose.png';
import eh_lose from '../assets/Lost.png';
import zh_draw from '../assets/draw.png';
import en_draw from '../assets/draw.jpg';
export default {
  components: {
    Popup
  },
  data() {
    return {
      imgList: [
        require('../assets/stone.png'),
        require('../assets/scissors.png'),
        require('../assets/cloth.png'),
      ],
      HistoryList:[],
      muted: true,
      opponentImgIndex: 2,
      ownImgIndex: 1,
      timer: null,
      win: 0,
      lose: 0,
      draw: 0,
      roomId: null,
      waitingShow: false,
      startGameShow: false,
      startGame: '',
      waiting: '',
      language: null,
      resultImg: en_win,
      resultShow: true,
    }
  },
  computed: {
    ...mapGetters(["msg"]),
    msg() {
      return store.getters.msg;
    }
  },
  watch: {
    muted(val) {
      if (!val) {
        this.$refs.music.play();
      } else {
        this.$refs.music.pause()
      }
    },
    msg(val) {
      this.global_callback(val)
    },
  },
  created() {
    this.language = getDefaultLang()
    this.roomId = Cookies.get('roomId');
    // this.startTime();
  },
  methods: {
    formatterTime(type) {
      const t = new Date()
      return formatterDate(t, type)
    },
    continueFun() {
      this.opponentImgIndex = null
      this.ownImgIndex = null    
      this.resultImg = null
      this.resultShow = false
      let params = {
        Cmd: 'ready',
        Room: this.roomId,
        Uuid: String(Cookies.get('Uuid'))
      }
      sendSock(params);
    },
    exit() {
      this.$router.replace({ path: '/' });
    },
    startTime() {
      let num = 0;
      const _this = this;
      this.timer = setInterval(function () {
        num++;
        if (num > 4) {
          num = 1;
        }
        let str = ''
        for (let i = 0; i < num; i++) {
          str = str + '.'
        }
        if (_this.startGameShow) {
          _this.startGame = _this.$t('startGame') + str
        } else if (_this.waitingShow) {
          _this.waiting = _this.$t('waiting') + str
        }
      }, 500)
    },
    onClick(val) {
      let params = {
        Cmd: 'guess',
        HandNum: String(val),
        Room: this.roomId,
        Uuid: String(Cookies.get('Uuid'))
      }
      sendSock(params);
      this.startTime()
      this.waitingShow = true
    },
    global_callback(msg) {
      if (msg.Cmd == 'full') {
        // 开始游戏了
        this.startGameShow = false
        clearInterval(this.timer)
      } else if (msg.Cmd == 'result') {
        let Uuid = String(Cookies.get('Uuid'))
        console.log(JSON.parse(msg.Data));
        let strObj = {
          time:this.formatterTime('yyyy-MM-dd hh:mm:ss'),
        }
        JSON.parse(msg.Data).forEach(item => {
          let e = JSON.parse(item)
          if (e.Uuid == Uuid) {
            this.ownImgIndex = Number(e.HandNum) - 1
            if(this.language == "zh"){
              strObj.youHandNum = Number(e.HandNum) == 1 || Number(e.HandNum) == '1' ? "你出的是石头" : (Number(e.HandNum) == 2 || Number(e.HandNum) == '2' ? "你出的是剪刀" : "你出的是布")
            }else{
              strObj.youHandNum = Number(e.HandNum) == 1 || Number(e.HandNum) == '1' ? "you -----> stone" : (Number(e.HandNum) == 2 || Number(e.HandNum) == '2' ? "you -----> scissors" : "you -----> cloth")
            }
            this.win = Number(e.WinCount)
            this.lose = Number(e.LoseCount)
            if (e.Win == 0 || e.Win == '0') {// 输了 
              if (this.language == 'zh') {
                this.resultImg = zh_lose
                strObj.isWin = "你输了"
              } else {
                this.resultImg = eh_lose
                strObj.isWin = "You lost"
              }
            } else if (e.Win == 1 || e.Win == '1') {// 赢了
              if (this.language == 'zh') {
                this.resultImg = zh_win
                strObj.isWin = "你赢了"
              } else {
                this.resultImg = en_win
                strObj.isWin = "You Win"
              }
            } else if (e.Win == 2 || e.Win == '2') {// 平局
              if (this.language == 'zh') {
                this.resultImg = zh_draw
                strObj.isWin = "平局"
              } else {
                this.resultImg = en_draw
                strObj.isWin = "draw"
              }
            }
          } else {
            this.opponentImgIndex = Number(e.HandNum) - 1
            if(this.language == "zh"){
              strObj.opponentHandNum = Number(e.HandNum) == 1 || Number(e.HandNum) == '1' ? "对手出的是石头" : (Number(e.HandNum) == 2 || Number(e.HandNum) == '2' ? "对手出的是剪刀" : "对手出的是布")
            }else{
              strObj.opponentHandNum = Number(e.HandNum) == 1 || Number(e.HandNum) == '1' ? "opponent -----> stone" : (Number(e.HandNum) == 2 || Number(e.HandNum) == '2' ? "opponent -----> scissors" : "opponent -----> cloth")
            }
          }
        })
        this.HistoryList.push(strObj)
        this.waitingShow = false
        clearInterval(this.timer)
        this.resultShow = true
      }
    },
  },
}
</script>

<style lang="less" scoped>
.game {
  width: 100%;
  height: 100vh;
  background-color: #36b3d5;

  .head_title {
    width: 100%;
    height: 10vh;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .main {
    display: flex;
    .statistics {
      width: 15%;
      .blico {
        font-size: 1.5rem;
        font-weight: 700;
        height: 5vh;
        display: flex;
        justify-content: center;
      }

      .win {
        color: #00fff2;

        span {
          margin-left: 0.625rem;
          margin-right: 0.625rem;
        }
      }

      .draw {
        color: #00ff15;

        span {
          margin-left: 0.625rem;
          margin-right: 0.625rem;
        }
      }

      .lose {
        color: red;

        span {
          margin-left: 0.625rem;
          margin-right: 0.625rem;
        }
      }
    }

    .competition {
      width: 70%;
      height: 40vh;
      display: flex;

      .block {
        width: 33%;
        display: flex;
        justify-content: center;
        align-items: center;
      }

      .own_venue,
      .opponent_venue {
        .frame {
          width: 60%;
          margin: 0 auto;

          .fist_area {
            width: 80%;
            height: 30vh;
            margin: 0 auto;
            display: flex;
            justify-content: center;
            align-items: center;
            background-color: #fff;
            border: 0.625rem solid #000;

            .img {
              width: 80%;
              height: 25vh;
            }

            .text {
              font-size: 5rem;
              font-weight: 700;
            }
          }

          .name {
            display: flex;
            justify-content: center;
            align-items: center;
            width: 100%;
            height: 8vh;
            font-size: 1.5rem;
            font-weight: 700;
          }
        }
      }

      .VS {
        font-size: 7.5rem;
        margin: 0 1.5rem;
        color: red;
        font-style: oblique;
      }
    }
    .fill_blcok{
      width: 15%;
    }
  }


  .skill_venue {
    width: 95%;
    height: 25vh;
    margin: 0 auto;
    background-color: #fff;
    margin-top: 1.875rem;
    border: 0.625rem solid #000;
    display: flex;
    justify-content: center;
    align-items: center;

    .skill0,
    .skill1,
    .skill2 {
      width: 10%;
      height: 18vh;
      margin-left: 1.25rem;
      display: flex;
      align-items: center;
      justify-content: center;
      border: 0.625rem solid #000;

      .img {
        width: 100%;
        height: 18vh;
      }
    }

    .skill0 {
      background-color: aqua;
    }

    .skill1 {
      background-color: #E0808B;
    }

    .skill2 {
      background-color: #84BC7F;
    }
  }
  .History{
    width: 50%;
    height: 19.5vh;
    color: #fff;
    background-color: #000;
    margin: 0 auto;
    overflow: auto;
    .title{
      font-size: 24px;
      font-weight: 600;
    }
  }


  .Popup {
    font-size: 4.25rem;
    font-weight: 700;
    white-space: 0.125rem;

    .btn {
      width: 12.5rem;
      height: 4.375rem;
      font-size: 1.75rem;
      font-weight: 700;
      white-space: 0.125rem;
      color: #fff;
      background: -webkit-linear-gradient(left, pink, skyblue);
      border: none;
    }
  }

  .result_Popup {
    width: 28.125rem;
    height: 19.375rem;
    background-color: #fff;

    .img_block {
      width: 100%;
      height: 15.625rem;

      .img {
        width: 100%;
        height: 15.625rem;
      }
    }

    .btnList {
      width: 100%;
      height: 3.75rem;
      display: flex;

      .btn {
        width: 50%;
        height: 3.75rem;
        background-image: linear-gradient(110.6deg, #b39ddb 7%, #969fde 47.7%, #18ffff 100.6%);
        border: none;
        margin: 0; 
      }
    }
  }
}
</style>