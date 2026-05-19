import axios from 'axios'

const API = axios.create({
    baseURL: 'https://api.pubgstars.com/v1'
});

export default {
    getGames: params => API.get('/games', params),
    getLeaderboard: params => API.get('games/leaderboard', params),
    getGamesHistory: params => API.get('games/history', params),

    getUser: params => API.get('/user', params),
    getUserGames: params => API.get('/user/games', params),
    getTransactionLog: params => API.get('/user/transactionlog', params),
    getWithdraw: (params, options) => API.post('/user/withdraw', params, options),
    getDeposit: (params, options) => API.post('/user/depositmoney', params, options),

    getGamePassword: (game, options) => API.post('/games/password', game, options),
    getGameUsers: (game, options) => API.post('/games/users', game, options),
    registerGame: (game, options) => API.post('games/register', game, options),
    unregisterGame: (game, options) => API.post('games/unregister', game, options),

    sendMessage: (params, options) => API.post('user/sendmessage', params, options)
}
