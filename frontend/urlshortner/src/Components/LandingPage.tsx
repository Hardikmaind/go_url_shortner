const posts = [
    {
      id: 1,
      title: 'Outstanding Service',
      href: '#',
      description:
        'this is the best service i have ever used. I am very happy with the service. I will recommend this to all',
      date: 'Nov 16, 2024',
      datetime: '2020-03-16',
      category: { title: 'Engineering', href: '#' },
      author: {
        name: 'Michael Foster',
        role: 'Practo',
        href: '#',
        imageUrl:
          'https://images.unsplash.com/photo-1519244703995-f4e0f30006d5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
      },
    },{
        id: 2,
        title: ' Good Service',
        href: '#',
        description:
          'best url shortner out there. This is even free',
        date: 'Mar 16, 2024',
        datetime: '2020-03-16',
        category: { title: 'Engineering', href: '#' },
        author: {
          name: 'Shravani Akude',
          role: 'VI',
          href: '#',
          imageUrl:
            'https://images.unsplash.com/photo-1519244703995-f4e0f30006d5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
        },
      },{
        id: 3,
        title: 'Very Impressed',
        href: '#',
        description:
          'It was the no login part that got me. I am very impressed with the service',
        date: 'Dec 16, 2024',
        datetime: '2020-03-16',
        category: { title: 'Student', href: '#' },
        author: {
          name: 'Sammy',
          role: 'Student',
          href: '#',
          imageUrl:
            'https://images.unsplash.com/photo-1519244703995-f4e0f30006d5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=facearea&facepad=2&w=256&h=256&q=80',
        },
      },
    // More posts...
  ]
  
  export default function LandingPage():JSX.Element {
    return (
      <div className="bg-gradient-to-tr  from-white-800 via-purple-300 to-green-400  mx-5 lg:mx-0 sm:py-10 rounded-lg">
        <div className="mx-auto max-w-7xl  px-6 lg:px-8">
          <div className="mx-auto max-w-2xl  lg:mx-0">
            <h2 className="text-pretty  font-semibold tracking-tight text-gray-900 sm:text-3xl">About ShortKing</h2>
            <p className=" text-lg/8 text-purple-700 font-bold">Shorten URLs for Free, No Login, Unlimited Convenience!</p>
          </div>
          <div className="mx-auto mt-10 grid max-w-2xl grid-cols-1 gap-x-8 gap-y-16 border-t border-gray-200 pt-10 sm:mt-2 lg:mx-0 lg:max-w-none lg:grid-cols-3">
            {posts.map((post) => (
              <article key={post.id} className="flex max-w-xl flex-col items-start justify-between">
                <div className="flex items-center gap-x-4 text-xs">
                  <time dateTime={post.datetime} className="text-white font-bold">
                    {post.date}
                  </time>
                  <a
                    href={post.category.href}
                    className="relative z-10 rounded-full bg-gray-50 px-3 py-1.5 font-medium text-gray-600 hover:bg-gray-100"
                  >
                    {post.category.title}
                  </a>
                </div>
                <div className="group relative">
                  <h3 className="mt-3  font-bold text-black text-xl group-hover:text-gray-600">
                    <a href={post.href}>
                      <span className="absolute inset-0" />
                      {post.title}
                    </a>
                  </h3>
                  <p className="mt-5 line-clamp-3 text-sm/6 text-white">{post.description}</p>
                </div>
                <div className="relative mt-8 flex items-center gap-x-4">
                  <img alt="" src={post.author.imageUrl} className="size-10 rounded-full bg-gray-50" />
                  <div className="text-sm/6">
                    <p className="font-semibold text-black text-xl">
                      <a href={post.author.href}>
                        <span className="absolute inset-0 " />
                        {post.author.name}
                      </a>
                    </p>
                    <p className="text-white underline  underline-offset-3">{post.author.role}</p>
                  </div>
                </div>
              </article>
            ))}
          </div>
        </div>
      </div>
    )
  }
  