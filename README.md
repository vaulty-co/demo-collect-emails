# An email list without emails

This is the source code of the demo application of the [Vaulty cookbook](https://docs.vaulty.co/cookbooks/email-list-with-protected-emails).

![front-page](https://docs.vaulty.co/img/cookbooks/demo-front.png)

Instructions on how to run the app, how to add Vaulty can be found in the [cookbook](https://docs.vaulty.co/cookbooks/email-list-with-protected-emails).

### Run the app

For the demo ap, we use Docker, Docker Compose. We also need the API key of the active Mailgun account. First, let's create `.env` file with environment variables:

```bash
MG_API_KEY=key-xxxxxxxxxxxxxxxxxxxxxx
MG_DOMAIN=mg.yourdomain.com
```

* `MG_DOMAIN` - Your Mailgun domain
* `MG_API_KEY` - Mailgun [API Key](https://help.mailgun.com/hc/en-us/articles/203380100-Where-Can-I-Find-My-API-Key-and-SMTP-Credentials-)

To run the demo application you need to put these commands into your shell:

```shell
git clone git@github.com:vaulty-co/demo-collect-emails.git
cd demo-collect-emails
docker-compose up
```

Then, navigate the browser to http://127.0.0.1:3000
