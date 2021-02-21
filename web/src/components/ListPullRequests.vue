<template>
  <div v-if="$store.state.token !== ''">
    <br/>
    <b-container>
      <div>
        <h4>My Pull Requests</h4>
      </div>
      <br/>
      <div>
        <b-table
            :busy="isBusy"
            :items="myPullRequests"
            :fields="myPullRequestsFields"
            striped
            responsive="sm"
            head-variant="light"
            small
        >
          <template #table-busy>
            <div class="text-center text-danger my-2">
              <b-spinner class="align-middle"></b-spinner>
              <strong>Loading...</strong>
            </div>
          </template>
          <template v-slot:cell(created_at)="row">
            <b :style="getDateStyle(row.item.created_at)">{{getFormattedDate(row.item.created_at)}}</b>
          </template>
          <template v-slot:cell(state)="row">
            <b v-if="row.item.state === 'COMMENTED'" style="color: #436f8a;">{{row.item.state}}</b>
            <b v-if="row.item.state === 'APPROVED'" style="color: #96bb7c;">{{row.item.state}}</b>
            <b v-if="row.item.state === 'PENDING'" style="color: #ffbd69;">{{row.item.state}}</b>
            <b v-if="row.item.state === 'CHANGES_REQUESTED'" style="color: red;">{{row.item.state}}</b>
            <b v-if="row.item.state === 'DISMISSED'" style="color: slateblue;">{{row.item.state}}</b>
          </template>
          <template v-slot:cell(link)="row">
            <a v-bind:href="row.item.url" target="”_blank”">Github</a>
          </template>
          <template v-slot:cell(link)="row">
            <a v-bind:href="row.item.url" target="”_blank”">Github</a>
          </template>
          <template v-slot:cell(story)="row">
            <a v-if="row.item.story !== ''" v-bind:href="getStory(row.item.story)" target="”_blank”">{{row.item.story}}</a>
          </template>
          <template #cell(details)="row" class="align-middle">
            <b-icon icon="plus-circle" @click="row.toggleDetails" v-if="!row.detailsShowing"></b-icon>
            <b-icon icon="dash-circle" @click="row.toggleDetails" v-if="row.detailsShowing"></b-icon>
          </template>

          <template #row-details="row">
            <b-card>
              <b-row
                  class="mb-2"
                  v-for="reviewRequest in row.item.review_requests"
                  :key="reviewRequest.requested_reviewer"
              >
                <b-col sm="3" class="text-sm-right">
                  <b>{{reviewRequest.requested_reviewer}}:</b>
                </b-col>
                <b-col>
                  <b v-if="reviewRequest.state === 'COMMENTED'" style="color: #436f8a;">{{reviewRequest.state}}</b>
                  <b v-if="reviewRequest.state === 'APPROVED'" style="color: #96bb7c;">{{reviewRequest.state}}</b>
                  <b v-if="reviewRequest.state === 'PENDING'" style="color: #ffbd69;">{{reviewRequest.state}}</b>
                  <b v-if="reviewRequest.state === 'CHANGES_REQUESTED'" style="color: red;">{{reviewRequest.state}}</b>
                  <b v-if="reviewRequest.state === 'DISMISSED'" style="color: slateblue;">{{reviewRequest.state}}</b>
                </b-col>
              </b-row>
            </b-card>
          </template>
          <template v-slot:custom-foot="data">
            <b-tr>
              <b-td colspan="7" variant="light" class="text-right">
                Total Rows: <b>{{data.items.length}}</b>
              </b-td>
            </b-tr>
          </template>
        </b-table>
      </div>
      <br>

      <div>
        <h4>Pull Requests to review</h4>
      </div>
      <br/>
      <div>
        <b-table
            :busy="isBusy"
            :items="pullRequestsToReview"
            :fields="pullRequestsToReviewFields"
            striped
            responsive="sm"
            head-variant="light"
            small
        >
          <template #table-busy>
            <div class="text-center text-danger my-2">
              <b-spinner class="align-middle"></b-spinner>
              <strong>Loading...</strong>
            </div>
          </template>
          <template v-slot:cell(created_at)="row">
            <b :style="getDateStyle(row.item.created_at)">{{getFormattedDate(row.item.created_at)}}</b>
          </template>
          <template v-slot:cell(state)="row">
            <b v-if="row.item.state === 'COMMENTED'" style="color: #436f8a;">{{row.item.state}}</b>
            <b v-if="row.item.state === 'APPROVED'" style="color: #96bb7c;">{{row.item.state}}</b>
            <b v-if="row.item.state === 'PENDING'" style="color: #ffbd69;">{{row.item.state}}</b>
            <b v-if="row.item.state === 'CHANGES_REQUESTED'" style="color: red;">{{row.item.state}}</b>
            <b v-if="row.item.state === 'DISMISSED'" style="color: slateblue;">{{row.item.state}}</b>
          </template>
          <template v-slot:cell(story)="row">
            <a v-if="row.item.story !== ''" v-bind:href="getStory(row.item.story)" target="”_blank”">{{row.item.story}}</a>
          </template>
          <template v-slot:cell(link)="row">
            <a v-bind:href="row.item.url" target="”_blank”">Github</a>

          </template>
          <template v-slot:custom-foot="data">
            <b-tr>
              <b-td colspan="7" variant="light" class="text-right">
                Total Rows: <b>{{data.items.length}}</b>
              </b-td>
            </b-tr>
          </template>
        </b-table>
      </div>
      <br/>

      <div>
        <h4>Reviews per User</h4>
      </div>
      <br/>
      <div>
        <b-table
            :busy="isBusy"
            :items="reviewersCount"
            :fields="reviewersCountFields"
            striped
            responsive="sm"
            head-variant="light"
            small
        >
        </b-table>
      </div>
    </b-container>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "ListPullRequests",
  props: {},
  data() {
    return {
      searchInterval:'',
      isBusy: true,
      myPullRequestsFields: [
        "created_at",
        "state",
        "application",
        "story",
        "link",
        "details",
      ],
      pullRequestsToReviewFields: [
        "created_at",
        "author",
        "state",
        "application",
        "story",
        "link",
      ],
      reviewersCountFields: [
        "username",
        "count",
      ],
      myPullRequests: [],
      pullRequestsToReview: [],
      reviewersCount: [],
    };
  },
  methods: {
    search() {
      axios
          .get("http://localhost:8081/pr-viewer/pull-requests", {
            params: this.queryParams(),
          })
          .then((res) => {
            this.myPullRequests = [];
            this.pullRequestsToReview = [];

            res.data.pull_requests.forEach(pr => {
              if (this.$store.state.username === pr.author) {
                this.myPullRequests.push(pr)
              } else {
                pr.review_requests.forEach(r => {
                  if (this.$store.state.username === r.requested_reviewer) {
                    pr.state = r.state
                  }
                });

                this.pullRequestsToReview.push(pr)
              }
            });

            this.reviewersCount = res.data.reviewers_count
            this.isBusy = false;
          })
          .catch((err) => {
            console.log(err);
          });
    },
    getDateStyle(date) {
      let DAY = 1000 * 60 * 60 * 24
      let now = new Date()
      let dateToParse = new Date(date)

      console.log(now - dateToParse)

      if ((now - dateToParse) < (2 * DAY)) {
        return "color: #96bb7c;"
      } else if ((now - dateToParse) < (4 * DAY)) {
        return "color: #ffbd69;"
      } else {
        return "color: red;"
      }
    },
    getFormattedDate(date) {
      const dateFormat = require('dateformat');

      return dateFormat(date, "yyyy, mmm dd");
    },
    getStory(number) {
      return "https://mercadolibre.atlassian.net/browse/" + number
    },
    queryParams() {
      return {
        token: this.$store.state.token,
      };
    },
  },
  created() {
    if (this.$store.state.token !== "") {
      this.search();
      //this.interval = setInterval(() => this.search(), 1000);
    }
    this.searchInterval = setInterval(() => {
      console.log ("llamando a search");
      //this.search();
    }, 1000 * 60 *5); // every 5 minutes
  },
  destroyed() {
    clearInterval(this.searchInterval)
  },
  watch: {
    "$store.state.token": function (val, oldVal) {
      console.log(val + oldVal);
      this.search();
    },
  },
};
</script>