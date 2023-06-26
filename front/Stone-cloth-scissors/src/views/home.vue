<template>
  <div class="home">
    <div class="main">
      <div class="game_name">
        {{ $t('gameName') }}
      </div>
      <div class="hang_block">
        <el-button class="btn" type="primary" @click="createOrjoinPopup('createRoom')">{{ $t('createRoom') }}</el-button>
        <el-button class="btn" type="primary" @click="createOrjoinPopup('joinRoom')">{{ $t('joinRoom') }}</el-button>
      </div>
    </div>
    <div class="rule">
      <h1><b>{{ $t('rule') }}</b></h1>
      <div class="rule_content">
        {{ $t('ruleContent') }}
      </div>
    </div>
    <Popup :show="maskShow">
      <div class="Popup">
        <div class="top">
          <div class="block" @click="closeMask">
            <i class="el-icon-close"></i>
          </div>
        </div>
        <div class="head">
          <h1><b>{{ popupTitle }}</b></h1>
        </div>
        <div class="main">
          <el-input class="input" v-model="roomId" :placeholder="$t('popupPlaceholder')"></el-input>
        </div>
        <div class="bottom">
          <el-button class="btn" type="primary" @click="createOrjoinRoom">{{ createOrjoinText }}</el-button>
        </div>
      </div>
    </Popup>
  </div>
</template>

<script>
import { sendSock } from "../utils/socket";
import { getDefaultLang } from "../utils/utils";
import Cookies from 'js-cookie'
import Popup from '@/components/mask/mask.vue';
export default {
  components:{
    Popup
  },
  data() {
    return {
      maskShow:false,
      roomId:'',
      language:'',
      popupTitle:'',
      createOrjoinText:'',
      timer:null,
    }
  },
  created() {
    this.language = getDefaultLang()
  },
  methods: {
    createOrjoinRoom(){
      let Uuid = new Date().getTime();
      Cookies.set("Uuid",Uuid, { expires: 1 });
      let params = {
        Room:this.roomId,
        Cmd:'login',
        User:'',
        AvatarUrl:'',
        Content:'',
        Uuid:String(Uuid),
        HandNum:'',
        GuessNum:'',
      }
      Cookies.set("roomId",this.roomId, { expires: 1 });
      sendSock(params);
    },
    createOrjoinPopup(val){
      this.createOrjoinText = this.language == 'zh' ? (val == 'createRoom' ? this.$t('create') : this.$t('Join')) : this.$t('createOrJoin');
      this.popupTitle = this.$t(val)
      this.maskShow = true
    },
    closeMask(){
      this.maskShow = false
    },
    // global_callback(msg) {
    //   if (msg.Cmd == 'loginsucess') {
    //     Cookies.set("roomId",this.roomId, { expires: 1 });
    //     this.$router.replace({path:'/game'});
    //   }
    // },
  },
}
</script>
<style lang="less" scoped>
.home {
  width: 100%;
  height: 100vh;
  background-image: linear-gradient(106.5deg, #ffd7b9e8 23%, #df9ff7cc 93%);

  .main {
    width: 100%;
    height: 45vh;
    .game_name{
      width: 100%;
      height: 10vh;
      display: flex;
      justify-content: center;
      align-items: center;
      font-size: 3rem;
      font-weight: 700;
    }
    .hang_block {
      margin: 0 auto;
      width: 40%;
      height: 35vh;
      display: flex;
      justify-content: space-around;
      align-items: center;

      .btn {
        min-width: 200px;
        min-height: 80px;
        font-size: 24px;
        font-weight: 600;
        background-image: linear-gradient(110.6deg, #b39ddb 7%, #969fde 47.7%, #18ffff 100.6%);
        border: none;
        &:hover{
          min-width: 240px;
          min-height: 100px;
          font-size: 28px;
        }
      }
    }
  }
  .rule{
    h1{
      font-size: 3rem;
      display: flex;
      justify-content: center;
      align-items: center;
      margin-bottom: 2rem;
      font-weight: 800;
    }
    .rule_content{
      width: 65%;
      font-weight: 600;
      font-size: 2rem;
      margin: 0 auto;
      text-indent: 3rem;
    }
  }
  .Popup{
    width: 31.25rem;
    height: 25rem;
    background-color: #fff;
    border: 0.625rem solid #000;
    position: relative;
    .top{
      width: 100%;height: 1.875rem;
      display: flex;
      justify-content: flex-end;
      .block{
        border-left: 0.625rem solid #000;
        border-bottom: 0.625rem solid #000;
        font-weight: 800;
        font-size: 1.25rem;
        display: flex;
        align-items: center;
        cursor: pointer;
        &:hover{
          background-color: #4cb8ea;
          color: #fff;
        }
      }
    }
    .head{
      width: 100%;height: 6.25rem;
      display: flex;
      justify-content: center;
      align-items: center;
    }
    .main{
      width: 100%;height: 7.5rem;
      display: flex;
      align-items: center;
      justify-content: center;
      .input{
        width: 70%;
        ::v-deep .el-input__inner{
          height: 3.75rem !important;
          font-size: 1.5rem;
        }
      }
    }
    .bottom{
      width: 100%;height: 9.375rem;
      display:flex;
      justify-content: center;
      align-items: center;
      .btn{
        width: 12.5rem;
        height: 3.75rem;
        font-size: 1.5rem;
        font-weight: 700;
      }
    }
  }
}
</style>