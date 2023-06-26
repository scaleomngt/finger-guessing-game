// The 'game_2ase7z' program.
program game_2ase7z.aleo {

    struct GameInfo {
        addr: address,
        gates: u64,
        game_id:  field,
        player_a: field,
        player_b: field
    }

    record PlayRecord {
        owner: address,
        gates: u64,
        game_id:  field,
        play_id:  field,
        player:   field,
        hand_num: u8
    }
    
    record GameResult {
        owner: address,
        gates: u64,
        game_id:  field,
        play_id:  field,
        player:   field,
        is_winner: u8 
    }
    
    transition create_game(
       public addr: address,
       public game_id: field, 
       public player:  field) -> GameInfo {

        return GameInfo {
          addr: addr,
          gates: 0u64,
          game_id: game_id,
          player_a: player,
          player_b: 0field,
        };
    }
    
    transition join_game(public gameInfo: GameInfo, public player: field) -> GameInfo {
        // Ensure the player_b is the empty.
        assert_eq(gameInfo.player_b, 0field);
        return GameInfo {
          addr: gameInfo.addr,
          gates: 0u64,
          game_id: gameInfo.game_id,
          player_a: gameInfo.player_a,
          player_b: player,
        };
    }

    
    //mapping plays: field => u8;
    
    transition start_game(gameInfo: GameInfo, play_id: field, player: field, handNum: u8) -> PlayRecord {
        assert_neq(gameInfo.player_b, 0field);
        
        return PlayRecord {
          owner: gameInfo.addr,
          gates: 0u64,
          game_id: gameInfo.game_id,
          play_id: play_id,
          player:  player,
          hand_num:handNum,
        };  
        //} then finalize(play_id, handNum);
        
    }
    //finalize start_game(play_id: field, handNum: u8) ->u8 {
    //    let current_value: u8 = Mapping::get_or_init(plays, play_id, 0u8);
    //    assert_eq(current_value, 0u8);
    //    Mapping::set(plays, play_id, handNum);
    //    return 0u8;
    //}
    
    //1 石头
    //2 剪刀
    //3 布
    transition finish_game(player_a: PlayRecord, player_b: PlayRecord  ) -> GameResult {
       assert_eq(player_a.game_id, player_b.game_id);
       assert_eq(player_a.play_id, player_b.play_id);
       assert_neq(player_a.player, player_b.player);
       assert_neq(player_a.hand_num, 0u8);
       assert_neq(player_b.hand_num, 0u8);
       
       if (player_a.hand_num > 3u8 ) {
          assert(false);
       }
       if (player_b.hand_num > 3u8 ) {
          assert(false);
       }
       
       let is_winner: u8 = get_game_result(player_a, player_b );
       
       return GameResult {
         owner: player_a.owner,
         gates: 0u64,
         game_id:  player_a.game_id,
         play_id:  player_a.play_id,
         player:   player_a.player,
         is_winner: is_winner, 
       };
    }
    
    
    function get_game_result(player_a: PlayRecord, player_b: PlayRecord) -> u8 {
       if (player_b.hand_num == player_a.hand_num) {
           return 0u8;
       }
       if (player_b.hand_num == 1u8) {
           if (player_a.hand_num == 2u8) {
             return 1u8;
           } else {
             return 2u8;
           }  
       } else if (player_b.hand_num == 2u8) {
           if (player_a.hand_num == 3u8) {
             return 1u8;
           } else {
             return 2u8;
           }
       } else if (player_b.hand_num == 3u8) {
           if (player_a.hand_num == 1u8) {
             return 1u8;
           } else {
             return 2u8;
           }
       }
       return 0u8;
    }
}
