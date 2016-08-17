# DAO Creator

Auto generate java class code from JSON URL.

it's a

- Use Gson
- Use AutoValue


## usage

```shell
~/g/s/g/s/dao-creator ❯❯❯ dao-creator -name "User" https://api.github.com/users/satoshun
public class User {
    @SerializedName("login") String login;
    @SerializedName("id") int id;
    @SerializedName("avatar_url") String avatarUrl;
    @SerializedName("gravatar_id") String gravatarId;
    @SerializedName("url") String url;
    @SerializedName("html_url") String htmlUrl;
    @SerializedName("followers_url") String followersUrl;
    @SerializedName("following_url") String followingUrl;
    @SerializedName("gists_url") String gistsUrl;
    @SerializedName("starred_url") String starredUrl;
    @SerializedName("subscriptions_url") String subscriptionsUrl;
    @SerializedName("organizations_url") String organizationsUrl;
    @SerializedName("repos_url") String reposUrl;
    @SerializedName("events_url") String eventsUrl;
    @SerializedName("received_events_url") String receivedEventsUrl;
    @SerializedName("type") String type;
    @SerializedName("site_admin") boolean siteAdmin;
    @SerializedName("name") String name;
    @SerializedName("company") String company;
    @SerializedName("blog") String blog;
    @SerializedName("location") String location;
    @SerializedName("email") String email;
    @SerializedName("hireable") boolean hireable;
    @SerializedName("bio") Object bio;
    @SerializedName("public_repos") int publicRepos;
    @SerializedName("public_gists") int publicGists;
    @SerializedName("followers") int followers;
    @SerializedName("following") int following;
    @SerializedName("created_at") Date createdAt;
    @SerializedName("updated_at") Date updatedAt;
}
```

or Use AutoValue.

```shell
@AutoValue public abstract class User {
    @SerializedName("login") abstract String login();
    @SerializedName("id") abstract int id();
    @SerializedName("avatar_url") abstract String avatarUrl();
    @SerializedName("gravatar_id") abstract String gravatarId();
    @SerializedName("url") abstract String url();
    @SerializedName("html_url") abstract String htmlUrl();
    @SerializedName("followers_url") abstract String followersUrl();
    @SerializedName("following_url") abstract String followingUrl();
    @SerializedName("gists_url") abstract String gistsUrl();
    @SerializedName("starred_url") abstract String starredUrl();
    @SerializedName("subscriptions_url") abstract String subscriptionsUrl();
    @SerializedName("organizations_url") abstract String organizationsUrl();
    @SerializedName("repos_url") abstract String reposUrl();
    @SerializedName("events_url") abstract String eventsUrl();
    @SerializedName("received_events_url") abstract String receivedEventsUrl();
    @SerializedName("type") abstract String type();
    @SerializedName("site_admin") abstract boolean siteAdmin();
    @SerializedName("name") abstract String name();
    @SerializedName("company") abstract String company();
    @SerializedName("blog") abstract String blog();
    @SerializedName("location") abstract String location();
    @SerializedName("email") abstract String email();
    @SerializedName("hireable") abstract boolean hireable();
    @SerializedName("bio") abstract Object bio();
    @SerializedName("public_repos") abstract int publicRepos();
    @SerializedName("public_gists") abstract int publicGists();
    @SerializedName("followers") abstract int followers();
    @SerializedName("following") abstract int following();
    @SerializedName("created_at") abstract Date createdAt();
    @SerializedName("updated_at") abstract Date updatedAt();
}
```


## install

Use Go compiler.

```
go install github.com/satoshun/dao-creator/cmd
```


## todos

- corresponds other library.
